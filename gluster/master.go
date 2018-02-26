package gluster

import (
	"fmt"
	"strconv"
	"math/rand"
	"strings"
	"encoding/gob"
	"path/filepath"
	"hash/crc32"
	"net"
	"reflect"
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
var file_list []common.FuncFileContent
var jobs_done []bool
var debugFlag = true;


/*
* Public Functions
*/
//execute function on remote node
//reply should be a pointer to the type returned by the function being called
//args is the arguments the function is expecting
//give back job id that can be used with wait or -1 on error
func RunDist(funct string, reply interface{}, args ...interface{}) int{
	//select the runner
	var cur_runner = pick_runner()
	if cur_runner == nil {
		return -1
	}

	//split string into function file and function
	var fun_elements = strings.Split(funct, ".")
	if(len(fun_elements) != 2){
		fmt.Println("Invalid functions string")
		return -1
	}

	//validity check (reply must be pointer)
	if(reflect.TypeOf(reply).Kind() != reflect.Ptr || reflect.TypeOf(reply).Elem().Kind() == reflect.Ptr){
		debugPrint("Return type is not a single pointer, aborting")
		return -1
	}

	//search imported files
	for _, el := range file_list {
		if(strings.Compare(el.File.CallPrefix, fun_elements[0]) == 0){
			//generate job id
			var id = len(jobs_done)
			jobs_done = append(jobs_done, false)

			//tell runner to run the function
			go runner_execute_function(cur_runner, id, fun_elements[1], el, reply, args)
			return id
		}
	}

	fmt.Println("Unable to locate function file: ", fun_elements[0])
	return -1
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
	var sum = crc32.ChecksumIEEE(contents)

	//add to list
	var new_func_file = common.FuncFileContent{}
	new_func_file.File.CallPrefix = call_name
	new_func_file.File.Checksum = sum
	new_func_file.FileType = common.GO_FILE
	new_func_file.Content = contents
	file_list = append(file_list, new_func_file)
}

//add a file of functions to be used
func ImportFunctionFileSO(filename string) {
	//must end in .go
	if(filepath.Ext(filename) != ".so"){
		fmt.Println("Invalid filename, must end in .so")
		return
	}

	//trim down filename to just the base without the extension
	var call_name = strings.TrimSuffix(filename, ".so")
	call_name = filepath.Base(call_name)

	//read whole file into memory
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to read file: ", err)
		return
	}

	//compute hash of file
	var sum = crc32.ChecksumIEEE(contents)

	//add to list
	var new_func_file = common.FuncFileContent{}
	new_func_file.File.CallPrefix = call_name
	new_func_file.File.Checksum = sum
	new_func_file.FileType = common.SO_FILE
	new_func_file.Content = contents
	file_list = append(file_list, new_func_file)
}

//returns whether the job with the given id is done executing
func JobDone(id int) bool {
	//check for invalid id
	if(id < 0 || id >= len(jobs_done)){
		return false
	}

	return jobs_done[id]
}

//Turn on debugging
func SetDebug(debug bool){
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

func compareType(want common.FuncSignature, haveRep interface{}, haveArgs []interface{}) bool {

	//check number of args
	if(len(want.In) != len(haveArgs)){
		debugPrint("Number of args does not match want", len(want.In), "have", len(haveArgs))
		return false
	}

	//compare args
	for i := 0; i < len(want.In); i++ {
		var haveType = common.EncodeType(reflect.TypeOf(haveArgs[i]))
		if(want.In[i] != haveType){
			debugPrint("Types don't match for arg", (i+1), "want:", want.In[i], ", have:", haveType)
			return false
		}
	}

	//compare return type
	var repType = common.EncodeType(reflect.TypeOf(haveRep).Elem())
	if(want.Out != repType){
		debugPrint("Types don't match for reply, want:", want.Out, ", have:", repType)
		return false
	}

	return true
}

func runner_execute_function(run *runner, id int, funct string, file common.FuncFileContent, reply interface{}, args []interface{}){
	conn, err := net.Dial("tcp", run.ip)
	if err != nil {
		fmt.Printf("Error connecting to runner")
		return
	}
	defer conn.Close()

	//setup gob
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	//send exec request
	var execReq = common.ExecRequest{file.File.Checksum, file.File.CallPrefix, funct}
	common.SendByte(conn, common.EXEC_FUNC)
	encoder.Encode(execReq)

	//get response
	var resp = common.RecvByte(conn)

	//function file needed, send it to runner
	if(resp == common.REQUESTING_FILE){
		common.SendByte(conn, common.FILE_INCOMING)
		encoder.Encode(file)
		resp = common.RecvByte(conn)
	} 
	
	
	var funcType common.FuncSignature
	//file is on runner, args needed
	if(resp == common.REQUESTING_ARGS){
		//read the function signature
		decoder.Decode(&funcType)
	} else{
		//TODO invalid response
		fmt.Println("Invalid response from runner")
		return
	}

	//TODO verify arguments match function signature
	debugPrint("Function type is", funcType)
	if(!compareType(funcType, reply, args)){
		debugPrint("Compare failed")
		return
	}


	//send all arguments
	common.SendByte(conn, common.ARGS_INCOMING)
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

	//make job id as done
	jobs_done[id] = true
}