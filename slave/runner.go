package main

import (
	"net"
	"fmt"
	"os/exec"
	"io/ioutil"
	"encoding/gob"
	"plugin"
	"reflect"
	"sync"
	"runtime"
	"../common"
)

/*
* globals
*/
var funcFileList []common.FuncFile
var funcListMut = &sync.RWMutex{}
var status common.RunnerStatus;

func main() {
	//TODO check go version and OS to make sure plugins can be built

	wait_for_work()
}


func wait_for_work() {
	//listen for connection from master
	listen , err := net.Listen("tcp", "localhost:1234");
	if err != nil {
		fmt.Println("Error setting up tcp connection: ", err);
		return;
	}
	defer listen.Close();

	//loop waiting for work connections
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err);
			return;
		}

		//thread work to be done
		go handle_work(conn);
	}
}

func handle_work(conn net.Conn) {

	var ctl = common.RecvByte(conn)
	debugPrint("Control received", ctl)

	//different cases
	if(ctl == common.EXEC_FUNC){
		exec_command(conn)
	} else if(ctl == common.GET_INFO){
		send_info(conn)
	}
}

func handle_hash_check(conn net.Conn) {
	debugPrint("Handling hash check")

	dec := gob.NewDecoder(conn)
	var hashRecv uint32
	var err = dec.Decode(&hashRecv)
	if(err != nil){
		fmt.Println("Error reading hash")
		return
	}

	debugPrint("Got hash: ", hashRecv)

	//check if hashes match
	var funcHere byte = common.NACK
	funcListMut.RLock()
	for _ , f := range funcFileList {
		if(f.Checksum == hashRecv){
			debugPrint("File is already here")
			funcHere = common.ACK
			break
		}
	}
	funcListMut.RUnlock()

	common.SendByte(conn, funcHere)

	handle_work(conn)
}

func recv_file(conn net.Conn) {
	debugPrint("Receiving file")

	dec := gob.NewDecoder(conn)

	file := &common.FuncFileContent{}
	dec.Decode(file)

	debugPrint("Got file with name: ", file.File.CallPrefix)

	//save to go file
	err := ioutil.WriteFile(file.File.CallPrefix + ".go", file.Content, 0644)
	if(err != nil){
		//TODO
		fmt.Println("Error writing go file")
	}

	//compile to a library
	cmd := exec.Command("go", "build", "-buildmode=plugin", 
		"-o", file.File.CallPrefix + string(file.File.Checksum) + ".so", 
		file.File.CallPrefix + ".go")
	err2 := cmd.Run()
	if(err2 != nil){
		fmt.Println("Failed to build")
		return
	}

	//add to list of available function files
	funcListMut.Lock()
	funcFileList = append(funcFileList, file.File)
	funcListMut.Unlock()

	//send ack
	common.SendByte(conn, common.ACK)

	handle_work(conn)
}

func check_for_file(conn net.Conn, execReq common.ExecRequest){
	//check if hashes match
	var funcHere = false
	funcListMut.RLock()
	for _ , f := range funcFileList {
		if(f.Checksum == execReq.Checksum){
			debugPrint("File is already here")
			funcHere = true
			break
		}
	}
	funcListMut.RUnlock()

	//if file is here, continue
	if(funcHere){
		return
	}

	//request file
	debugPrint("Requesting file")
	common.SendByte(conn, common.REQUESTING_FILE)

	var resp = common.RecvByte(conn)
	if(resp != common.FILE_INCOMING){
		//TODO error
	}
	
	//read file from network
	dec := gob.NewDecoder(conn)
	file := &common.FuncFileContent{}
	dec.Decode(file)

	debugPrint("Got file with name: ", file.File.CallPrefix)

	//save to go file
	err := ioutil.WriteFile(file.File.CallPrefix + ".go", file.Content, 0644)
	if(err != nil){
		//TODO
		fmt.Println("Error writing go file")
	}

	//compile to a library
	cmd := exec.Command("go", "build", "-buildmode=plugin", 
		"-o", file.File.CallPrefix + ".so", 
		file.File.CallPrefix + ".go")
	err2 := cmd.Run()
	if(err2 != nil){
		fmt.Println("Failed to build")
		return
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
	if(funcType.NumOut() > 0){
		funcSig.Out = common.EncodeType(funcType.Out(0))
	}

	return funcSig
}

func exec_command(conn net.Conn){
	debugPrint("Executing Function")

	//setup gob
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	//receive exec request
	exec := &common.ExecRequest{}
	decoder.Decode(exec)

	//make sure the function file is here
	check_for_file(conn, *exec)

	//load in function file library
	var path = exec.FuncFileName + ".so"
	p, err := plugin.Open(path)
	if(err != nil){
		//TODO send error response
		fmt.Println("Error opening function library: ", path, " : ", err)
		return
	}

	//lookup function
	f, err := p.Lookup(exec.FuncName)
	if(err != nil){
		//TODO send error response
		fmt.Println("Unable to find function in library: ", exec.FuncName, " : ", err)
		return 
	}

	//get type of function to be called
	var funcType = reflect.TypeOf(f)

	//request args and send function signature
	debugPrint("Requesting Args")
	common.SendByte(conn, common.REQUESTING_ARGS)
	encoder.Encode(encodeFuncSig(funcType))


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


	//call function
	var reply = reflect.ValueOf(f).Call(args)

	//send ack
	common.SendByte(conn, common.ACK)

	//encode and send back reply
	if(len(reply) > 0){
			sendReply(conn, reply[0].Interface())
	}

	debugPrint("Done calling function")
}

func sendReply(conn net.Conn, reply interface{}){
	//if response is pointer, dereference
	debugPrint("Sending reply", reply)
	enc := gob.NewEncoder(conn)
	enc.Encode(reply)
}


func send_info(conn net.Conn){
	var info = common.RunnerInfo{}
	info.Cores = runtime.NumCPU()
	info.Arch = runtime.GOARCH
	info.OS = runtime.GOOS

	//write structure
	sendReply(conn, info)
}

/*
* Debug functions
*/
var debugFlag = true
func debugPrint(args ...interface{}){
	if(debugFlag){
		fmt.Println(args)
	}
}