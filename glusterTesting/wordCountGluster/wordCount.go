package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
	"strings"
	"io/ioutil"
	"sync"
	"../../gluster/src/master"
)

func mergeMaps(mapList []map[string]int) map[string]int{
	var wg sync.WaitGroup
	var resultMap = struct{
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}
	for _, i := range mapList {
		go func(subMap map[string]int) {
			defer wg.Done()
			for word, count := range subMap {
				resultMap.RLock()
				gcount, ok := resultMap.m[word]
				resultMap.RUnlock()
				if ok {
					resultMap.Lock()
					resultMap.m[word] = gcount + count
					resultMap.Unlock()
				} else {
					resultMap.Lock()
					resultMap.m[word] = count
					resultMap.Unlock()
				}
			}
		} (i)
	}
	wg.Wait()
	return resultMap.m
}

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


	// Initalize Dataset into array of words "words"
	timer0 := time.Now()
	dat, err := ioutil.ReadFile("words.txt")
	check(err)
	book := strings.Replace(string(dat), "—", " ", -1)
	book = book + " "
	book = strings.Repeat(book, testSize)
	words := strings.Fields(book)
	timer1 := time.Now()

	//Initialize Gluster
	addRunners() // - This should be added to the library
	if debug == 1 {
		gluster.SetDebug(true)
	} else {
		gluster.SetDebug(false)
	}
	gluster.ImportFunctionFile("functions/functions.go")
	timer2 := time.Now()

	//Launch Runners
	var output map[string]int

	var runnerList []int
	for runnerID := 0; runnerID < runnerCount; runnerID++ {
		fmt.Println("Launching Runner: ", runnerID)
		var start = runnerID * len(words) / runnerCount
		var end = (runnerID + 1) * len(words) / runnerCount
		runnerList = append(runnerList, gluster.RunDist("functions.wordCount", reflect.TypeOf(output), words[start:end], processCount))
	}

	//Collect runner outputs
	var partialResults []map[string]int
	for _, runner := range runnerList {
		for !(gluster.JobDone(runner)) {
		}
		var partialResult = gluster.GetReturn(runner).(map[string]int)
		partialResults = append(partialResults, partialResult)
	}
	resultMap := mergeMaps(partialResults)
	timer3 := time.Now()
	
	fmt.Println("length of map: ", len(resultMap))
	fmt.Println("Dataset Init Time: ", timer1.Sub(timer0))
	fmt.Println("Gluster Init Time: ", timer2.Sub(timer1))
	fmt.Println("wordCount Time:    ", timer3.Sub(timer2))
}
