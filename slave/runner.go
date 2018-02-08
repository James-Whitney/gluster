package main

import (
	"net"
	"fmt"
	"os/exec"
	"io/ioutil"
	"encoding/gob"
	"plugin"
	"time"
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

	buf := make([]byte, 1)
	var _, err = conn.Read(buf)
	if(err != nil){
		fmt.Println("Error reading response", err)
		return
	}

	//different cases
	if(buf[0] == common.HASH_CMP){
		handle_hash_check(conn)
	} else if(buf[0] == common.SEND_FILE){
		recv_file(conn)
	} else if(buf[0] == common.EXEC_FUNC){
		exec_command(conn)
	}
}

func handle_hash_check(conn net.Conn) {
	fmt.Println("Handling hash check")

	hashBuf := make([]byte, 32)
	var _, err = conn.Read(hashBuf)
	if(err != nil){
		fmt.Println("Error reading hash")
		return
	}

	fmt.Println("Got hash: ", hashBuf)

	sendBuf := make([]byte, 1)
	sendBuf[0] = 0
	conn.Write(sendBuf)

	handle_work(conn)
}

func recv_file(conn net.Conn) {
	fmt.Println("Receiving file")

	dec := gob.NewDecoder(conn)

	file := &common.FuncFile{}
	dec.Decode(file)

	fmt.Println("Got file with name: ", file.CallPrefix)

	//save to go file
	err := ioutil.WriteFile(file.CallPrefix + ".go", file.Contents, 0)
	if(err != nil){
		//TODO
	}

	//TODO hacky, need to wait for file to appear in os filesystem
	time.Sleep(10 * time.Second)

	//compile to a library
	cmd := exec.Command("go", "build", "-buildmode=plugin", file.CallPrefix + ".go")
	err2 := cmd.Run()
	if(err2 != nil){
		fmt.Println("Failed to build")
		return
	}

	handle_work(conn)
}

func exec_command(conn net.Conn){
	fmt.Println("Executing Function")

	//receive command
	dec := gob.NewDecoder(conn)
	exec := &common.ExecSend{}
	dec.Decode(exec)

	var path = exec.FuncFile + ".so"

	fmt.Println("looking for ", path)

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
	reflect.ValueOf(f).Call(args)

	fmt.Println("Done calling function")
}