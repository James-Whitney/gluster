package gluster

import (
	"net/rpc"
	"fmt"
	"strconv"
	"math/rand"
	//. "../functions"
)

/*
* Constants
*/
const default_rpc_port = 1234

/*
* Structures
*/
//info for runner on each slave machine
type runner struct {
	//TODO connection, uptime, etc
	conn *rpc.Client
}

/*
* Globals
*/
var runner_list []runner


/*
* Public Functions
*/
//execute give rpc function on a runner
func RunDist(funct string, args interface{}, reply interface{}){
	var cur_runner = pick_runner()
	if cur_runner == nil {
		return
	}
	fmt.Println("abc")
	var err = cur_runner.conn.Call(funct, args, reply)
	fmt.Println("cds")
	if err != nil {
		fmt.Println("RPC failed: ", err)
	}
}

//add runner at the given ip with the default port
func AddRunner(ip string) {
	if runner_list == nil {
		runner_list = make([]runner, 0)
		fmt.Println("Init runner list")
	}


	var slave, err = rpc.Dial("tcp", ip + ":" + strconv.Itoa(default_rpc_port))
	if err != nil {
		fmt.Println("Unable to add runner: ", err)
		return
	}
	var new_runner = runner{slave}
	runner_list = append(runner_list, new_runner)
}

//add runner at the given ip with the given port
func AddRunnerPort(ip string, port int){
	if runner_list == nil {
		runner_list = make([]runner, 0)
	}

	var slave, err = rpc.Dial("tcp", ip + ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("Unable to add runner: ", err)
		return
	}
	var new_runner = runner{slave}
	runner_list = append(runner_list, new_runner)
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

