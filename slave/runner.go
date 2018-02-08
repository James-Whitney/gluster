package main

import (
	"net"
	"fmt"
	"os/exec"
	"io/ioutil"
	"encoding/gob"
	"plugin"
	"reflect"
	"../common"
)

func main() {
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
	if(ctl == common.HASH_CMP){
		handle_hash_check(conn)
	} else if(ctl == common.SEND_FILE){
		recv_file(conn)
	} else if(ctl == common.EXEC_FUNC){
		exec_command(conn)
	}
}

func handle_hash_check(conn net.Conn) {
	debugPrint("Handling hash check")

	hashBuf := make([]byte, 32)
	var _, err = conn.Read(hashBuf)
	if(err != nil){
		fmt.Println("Error reading hash")
		return
	}

	debugPrint("Got hash: ", hashBuf)

	//TODO check if hashes match
	common.SendByte(conn, common.NACK)

	handle_work(conn)
}

func recv_file(conn net.Conn) {
	debugPrint("Receiving file")

	dec := gob.NewDecoder(conn)

	file := &common.FuncFile{}
	dec.Decode(file)

	debugPrint("Got file with name: ", file.CallPrefix)

	//save to go file
	err := ioutil.WriteFile(file.CallPrefix + ".go", file.Contents, 0)
	if(err != nil){
		//TODO
		fmt.Println("Error writing go file")
	}

	//TODO hacky, need to wait for file to appear in os filesystem
	//time.Sleep(10 * time.Second)

	//compile to a library
	cmd := exec.Command("go", "build", "-buildmode=plugin", file.CallPrefix + ".go")
	err2 := cmd.Run()
	if(err2 != nil){
		fmt.Println("Failed to build")
		return
	}

	//send ack
	common.SendByte(conn, common.ACK)

	handle_work(conn)
}

func exec_command(conn net.Conn){
	debugPrint("Executing Function")

	//receive command
	dec := gob.NewDecoder(conn)
	exec := &common.ExecSend{}
	dec.Decode(exec)

	var path = exec.FuncFile + ".so"

	debugPrint("looking for ", path)

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
	var args []reflect.Value
	//get each argument and decode it
	for i := 0; i < funcType.NumIn(); i++ {
		var tmpArg = reflect.New(funcType.In(i))
		err = dec.Decode(tmpArg.Interface())
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


/*
* Debug functions
*/
var debugFlag = true
func debugPrint(args ...interface{}){
	if(debugFlag){
		fmt.Println(args)
	}
}