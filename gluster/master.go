package gluster

import (
	"fmt"
	"strconv"
	"math/rand"
	"strings"
	"encoding/gob"
	"path/filepath"
	"crypto/sha256"
	"net"
	"io/ioutil"
	"../common"
)

/*
* Structures
*/
//info for runner on each slave machine
type runner struct {
	//TODO connection, uptime, etc
	conn net.Conn
	ip string
}


/*
* Globals
*/
var runner_list []runner
var file_list []common.FuncFile


/*
* Public Functions
*/
//execute give rpc function on a runner
//funct should be of form file.function where file is the base filename without the extention
func RunDist(funct string, args interface{}, reply interface{}){
	//select the runner
	var cur_runner = pick_runner()
	if cur_runner == nil {
		return
	}

	//split string into function file and function
	var fun_elements = strings.Split(funct, ".")
	if(len(fun_elements) != 2){
		fmt.Println("Invalid functions string")
		return
	}

	//search imported files
	for _, el := range file_list {
		if(strings.Compare(el.CallPrefix, fun_elements[0]) == 0){
			//tell runner to run the function
			//TODO verify function is in file on server side
			runner_execute_function(cur_runner, fun_elements[1], args, reply, el)
			return
		}
	}

	fmt.Println("Unable to locate function file: ", fun_elements[0])
}

//add runner at the given ip with the default port
func AddRunner(ip string) {
	if runner_list == nil {
		runner_list = make([]runner, 0)
	}


	var full_ip = ip + ":" + strconv.Itoa(common.DEFAULT_PORT) 
	var slave, err = net.Dial("tcp", full_ip)
	if err != nil {
		fmt.Println("Unable to add runner: ", err)
		return
	}
	var new_runner = runner{slave, full_ip}
	runner_list = append(runner_list, new_runner)
}

//add runner at the given ip with the given port
func AddRunnerPort(ip string, port int){
	if runner_list == nil {
		runner_list = make([]runner, 0)
	}

	var full_ip = ip + ":" + strconv.Itoa(port)
	var slave, err = net.Dial("tcp", full_ip)
	if err != nil {
		fmt.Println("Unable to add runner: ", err)
		return
	}
	var new_runner = runner{slave, full_ip}
	runner_list = append(runner_list, new_runner)
}

//add a file of functions to be used
func ImportFunctionFile(filename string) {
	//must end in .go
	if(filepath.Ext(filename) != ".go"){
		fmt.Println("Invalid filename, must end in .go")
		return
	}

	//trim down filename to just the base without the extension
	var call_name = strings.TrimSuffix(filename, ".go")
	call_name = filepath.Base(call_name)

	//read whole file into memory
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to read file: ", err)
		return
	}

	//compute hash of file
	h := sha256.New()
	h.Write(contents)
	var sum = h.Sum(nil)

	//add to list
	var new_func_file = common.FuncFile{call_name, contents, sum}
	file_list = append(file_list, new_func_file)
}

/*
* Private Functions
*/
func pick_runner() *runner{
	//pick random
	if len(runner_list) == 0 {
		fmt.Println("No runners available!!")
		return nil
	}
	return &runner_list[rand.Intn(len(runner_list))]
}

func runner_compare_hash(conn net.Conn, hash []byte) bool{
	buf := make([]byte, 1)
	
	buf[0] = common.HASH_CMP
	conn.Write(buf)
	var n, _ = conn.Write(hash)
	
	fmt.Println("Wrote hash bytes: ", n)

	var _, err = conn.Read(buf)
	if(err != nil){
		fmt.Println("Error reading response")
		return false
	}

	//hash does match, function is already there
	if(buf[0] == 1){
		return true
	}

	return false
}

func runner_send_file(conn net.Conn, file common.FuncFile){
	buf := make([]byte, 1)
	buf[0] = common.SEND_FILE	
	conn.Write(buf)
	
	//encode file info
	encoder := gob.NewEncoder(conn)
	encoder.Encode(file)
}

func runner_execute_function(run *runner, funct string, args interface{}, reply interface{}, file common.FuncFile){
	conn, err := net.Dial("tcp", run.ip)
	if err != nil {
		fmt.Printf("Error connecting to runner")
	}
	defer conn.Close()

	if(!runner_compare_hash(conn, file.Checksum)){
		runner_send_file(conn, file)
	}
	conn.Close()

	conn2, _ := net.Dial("tcp", run.ip)
	//send exec command
	buf := make([]byte, 1)
	buf[0] = common.EXEC_FUNC	
	conn2.Write(buf)

	gob.Register(args)
	var execSend = common.ExecSend{file.CallPrefix, funct, args}
	encoder := gob.NewEncoder(conn2)
	encoder.Encode(execSend)
	fmt.Println("Sent command")

	var resp = &common.ExecResp{}
	dec := gob.NewDecoder(conn2)
	dec.Decode(resp)

}