package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"../../gluster/src/master"
)


func check(e error) {
    if e != nil {
        panic(e)
    }
}

func addRunners() {
	file, err := os.Open("./slaves.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	str, _, err := reader.ReadLine()
	for err == nil {
		gluster.AddRunner(string(str))
		str, _, err = reader.ReadLine()
	}
}

func main() {
	testSize, _ := strconv.Atoi(os.Args[1])
	runnerCount, _ := strconv.Atoi(os.Args[2])
	processCount, _ := strconv.Atoi(os.Args[3])
	debug, _ := strconv.Atoi(os.Args[4])


	//Initialize Gluster
	addRunners() // - This should be added to the library
	if debug == 1 {
		gluster.SetDebug(true)
	} else {
		gluster.SetDebug(false)
	}
	gluster.ImportFunctionFile("functions/functions.go")

	var input map[string]int
	input["test"] = 100
	
	//Launch Runners
	var output map[string]int

	var runnerList []int
	for runnerID := 0; runnerID < runnerCount; runnerID++ {
		fmt.Println("Launching Runner: ", runnerID)
		runnerList = append(runnerList, gluster.RunDist("functions.ReturnMap", reflect.TypeOf(output), input))
	}

	//Collect runner outputs
	for _, runner := range runnerList {
		for !(gluster.JobDone(runner)) {
		}
		var result = gluster.GetReturn(runner).(map[string]int)
	}
	
	fmt.Println("length of map: ", len(result))
	fmt.Println("map[test]: ", result["test"])
}
