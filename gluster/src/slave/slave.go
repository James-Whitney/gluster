package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"plugin"
	"reflect"
	"runtime"
	"sync"

	"../common"
)

/*
* globals
 */
var funcFileList []common.FuncFile
var funcListMut = &sync.RWMutex{}
var status common.RunnerStatus

func main() {
	//TODO check go version and OS to make sure plugins can be built

	wait_for_work()
}

func wait_for_work() {
	//listen for connection from master
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error setting up tcp connection: ", err)
		return
	}
	defer listen.Close()

	//loop waiting for work connections
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			return
		}

		//thread work to be done
		go handle_work(conn)
	}
}

func handle_work(conn net.Conn) {

	var rw = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	var ctl = common.RecvByte(rw)
	debugPrint("Control received", ctl)

	//different cases
	if ctl == common.EXEC_FUNC {
		exec_command(rw)
	} else if ctl == common.GET_INFO {
		send_info(rw)
	}
}

func check_for_file(rw *bufio.ReadWriter, execReq common.ExecRequest) {
	//check if hashes match
	var funcHere = false
	funcListMut.RLock()
	for _, f := range funcFileList {
		if f.Checksum == execReq.Checksum {
			debugPrint("File is already here")
			funcHere = true
			break
		}
	}
	funcListMut.RUnlock()

	//if file is here, continue
	if funcHere {
		return
	}

	//request file
	debugPrint("Requesting file")
	common.SendByte(rw, common.REQUESTING_FILE)

	var resp = common.RecvByte(rw)
	if resp != common.FILE_INCOMING {
		//TODO error
	}

	//read file from network
	dec := gob.NewDecoder(rw)
	file := &common.FuncFileContent{}
	dec.Decode(file)

	//go file received
	if file.FileType == common.GO_FILE {
		debugPrint("Got Go file with name: ", file.File.CallPrefix)

		//save to go file
		err := ioutil.WriteFile(file.File.CallPrefix+".go", file.Content, 0644)
		if err != nil {
			//TODO
			fmt.Println("Error writing go file")
		}

		//compile to a library
		cmd := exec.Command("go", "build",
			"-ldflags", "\"-pluginpath=plugin/hot-"+fmt.Sprint(file.File.Checksum)+"\"",
			"-buildmode=plugin",
			"-o", file.File.CallPrefix+fmt.Sprint(file.File.Checksum)+".so",
			file.File.CallPrefix+".go")
		err2 := cmd.Run()
		if err2 != nil {
			fmt.Println("Failed to build")
			return
		}
	} else if file.FileType == common.SO_FILE {
		debugPrint("Got SO file with name: ", file.File.CallPrefix)

		//save to SO file
		err := ioutil.WriteFile(file.File.CallPrefix+fmt.Sprint(file.File.Checksum)+".so", file.Content, 0644)
		if err != nil {
			//TODO
			fmt.Println("Error writing SO file")
		}

	} else {
		//TODO error
	}

	//add to list of available function files
	funcListMut.Lock()
	funcFileList = append(funcFileList, file.File)
	funcListMut.Unlock()
}

func encodeFuncSig(funcType reflect.Type) common.FuncSignature {
	var funcSig = common.FuncSignature{}

	//string for args
	for i := 0; i < funcType.NumIn(); i++ {
		funcSig.In = append(funcSig.In, common.EncodeType(funcType.In(i)))
	}

	//string for return
	if funcType.NumOut() > 0 {
		funcSig.Out = common.EncodeType(funcType.Out(0))
	}

	return funcSig
}

func exec_command(rw *bufio.ReadWriter) {
	debugPrint("Executing Function")

	//setup gob
	encoder := gob.NewEncoder(rw)
	decoder := gob.NewDecoder(rw)

	//receive exec request
	exec := &common.ExecRequest{}
	decoder.Decode(exec)

	//make sure the function file is here
	check_for_file(rw, *exec)

	//load in function file library
	var path = exec.FuncFileName + fmt.Sprint(exec.Checksum) + ".so"
	debugPrint("Loading plugin", path)
	p, err := plugin.Open(path)
	if err != nil {
		//TODO send error response
		fmt.Println("Error opening function library: ", path, " : ", err)
		return
	}

	//lookup function
	f, err := p.Lookup(exec.FuncName)
	if err != nil {
		//TODO send error response
		fmt.Println("Unable to find function in library: ", exec.FuncName, " : ", err)
		return
	}

	//get type of function to be called
	var funcType = reflect.TypeOf(f)

	//request args and send function signature
	debugPrint("Requesting Args")
	common.SendByte(rw, common.REQUESTING_ARGS)
	encoder.Encode(encodeFuncSig(funcType))
	rw.Flush()

	if common.RecvByte(rw) != common.ARGS_INCOMING {
		return
	}

	//get each argument and decode it
	var args []reflect.Value
	for i := 0; i < funcType.NumIn(); i++ {
		var tmpArg = reflect.New(funcType.In(i))
		err = decoder.Decode(tmpArg.Interface())
		if err != nil {
			fmt.Println("Error decoding argument", i)
		}
		args = append(args, tmpArg.Elem())
	}

	//send ack
	common.SendByte(rw, common.ACK)

	debugPrint("Got all args, calling function")

	//call function
	var reply = reflect.ValueOf(f).Call(args)

	//encode and send back reply
	if len(reply) > 0 {
		// fmt.Println(len(reply[0].Interface().([]int)))
		sendReply(rw, reply[0].Interface()) 
	}

	debugPrint("Done calling function, response sent")
}

func sendReply(conn *bufio.ReadWriter, reply interface{}) {
	//if response is pointer, dereference
	debugPrint("Sending reply") //, reply)
	debugPrint(reflect.TypeOf(reply))

	enc := gob.NewEncoder(conn)
	err := enc.Encode(reply)
	if err != nil {
		fmt.Println("error encoding args", err)
	}
	conn.Flush()
}

func send_info(rw *bufio.ReadWriter) {
	var info = common.RunnerInfo{}
	info.Cores = runtime.NumCPU()
	info.Arch = runtime.GOARCH
	info.OS = runtime.GOOS

	//write structure
	sendReply(rw, info)
}

/*
* Debug functions
 */
var debugFlag = true

func debugPrint(args ...interface{}) {
	if debugFlag {
		fmt.Println(args)
	}
}
