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
	//"reflect"
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
var debugFlag = true;


/*
* Public Functions
*/
//execute function on remote node
//reply should be a pointer to the type returned by the function being called
//args is the arguments the function is expecting
func RunDist(funct string, reply interface{}, args ...interface{}){
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

	//TODO check validity (reply must be pointer), arguments/reply must match type

	//search imported files
	for _, el := range file_list {
		if(strings.Compare(el.CallPrefix, fun_elements[0]) == 0){
			//tell runner to run the function
			//TODO verify function is in file on server side
			runner_execute_function(cur_runner, fun_elements[1], el, reply, args)
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

	//TODO check the file builds and create list of possible functions to call

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

//Turn on debugging
func setDebug(debug bool){
	debugFlag = debug;
}

/*
* Private Functions
*/
func debugPrint(args ...interface{}){
	if(debugFlag){
		fmt.Println(args)
	}
}

func pick_runner() *runner{
	//pick random
	if len(runner_list) == 0 {
		fmt.Println("No runners available!!")
		return nil
	}
	return &runner_list[rand.Intn(len(runner_list))]
}

func runner_compare_hash(conn net.Conn, hash []byte) bool{
	//send control byte
	common.SendByte(conn, common.HASH_CMP)
	//send hash
	var n, _ = conn.Write(hash)
	
	debugPrint("Wrote hash bytes: ", n)

	return common.RecvACK(conn)
}

func runner_send_file(conn net.Conn, file common.FuncFile){
	//send control byte
	common.SendByte(conn, common.SEND_FILE)
	
	//encode file info
	encoder := gob.NewEncoder(conn)
	encoder.Encode(file)

	common.RecvACK(conn)
}

func runner_execute_function(run *runner, funct string, file common.FuncFile, reply interface{}, args []interface{}){
	conn, err := net.Dial("tcp", run.ip)
	if err != nil {
		fmt.Printf("Error connecting to runner")
		return
	}
	defer conn.Close()

	if(!runner_compare_hash(conn, file.Checksum)){
		runner_send_file(conn, file)
	}

	//send control byte
	debugPrint("Sending exec cmd")
	common.SendByte(conn, common.EXEC_FUNC)

	//send which function to call
	var execSend = common.ExecSend{file.CallPrefix, funct}
	encoder := gob.NewEncoder(conn)
	encoder.Encode(execSend)

	//send all arguments
	for _, arg := range args{
		encoder.Encode(arg)
	}

	debugPrint("Sent command")

	//receive ack
	if(!common.RecvACK(conn)){
		//got NACK
		return
	}

	//get back response
	if(reply != nil){
		dec := gob.NewDecoder(conn)
		//debugPrint(reflect.TypeOf(reply).Elem())
		//var tmpReply = reflect.New(reflect.TypeOf(reply).Elem())
		//var tmpReply int
		err = dec.Decode(reply)
		if err != nil {
        	fmt.Println("Error decoding reply")
		}
		//debugPrint(tmpReply)
	}
}