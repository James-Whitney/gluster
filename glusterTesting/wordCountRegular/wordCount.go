package main

import (
	"fmt"
	"strings"
	"time"
	"io/ioutil"
	"sync"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

// func addRunners() {
// 	file, err := os.Open("./slaves.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	reader := bufio.NewReader(file)
// 	str, _, err := reader.ReadLine()
// 	for err == nil {
// 		gluster.AddRunner(string(str))
// 		str, _, err = reader.ReadLine()
// 	}
// }

func main() {

	// Initalize Dataset into array of words "words"
	timer0 := time.Now()
	dat, err := ioutil.ReadFile("words.txt")
	check(err)
	book := strings.Replace(string(dat), "—", " ", -1)
	book = book + " "
	book = strings.Repeat(book, 10)
	words := strings.Fields(book)
	timer1 := time.Now()

	//Initialize Gluster
	// addRunners() // - This should be added to the library
	// if debug == 1 {
	// 	gluster.SetDebug(true)
	// } else {
	// 	gluster.SetDebug(false)
	// }
	// gluster.ImportFunctionFile("functions/functions.go")
	timer2 := time.Now()

	// perform wordCount
	processCount := 4 // number of goRoutines per machine

	var wg sync.WaitGroup
	var globalDict = struct{
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}

	// globalDict.dict = make(map[string]int)
	for p := 0; p < processCount; p++ {

		wg.Add(1)
		go func(processID int) {
			defer wg.Done()
			var start = processID * len(words) / processCount
			var end = (processID + 1) * len(words) / processCount
			
			// var m map[string]int
			m := make(map[string]int)

			for i := start; i < end; i++ {
				word := strings.ToLower(strings.Trim(words[i], "*!(),.?;“”’_"))
				count := m[word]
				if count == 0 {
					m[word] = 1
				} else {
					m[word] = count + 1
				}
			}

			for word, count := range m {
				globalDict.RLock()
				gcount, ok := globalDict.m[word]
				globalDict.RUnlock()
				if ok {
					globalDict.Lock()
					globalDict.m[word] = gcount + count
					globalDict.Unlock()
				} else {
					globalDict.Lock()
					globalDict.m[word] = count
					globalDict.Unlock()
				}
			}
		} (p)
	}

	wg.Wait()

	timer3 := time.Now()
	
	fmt.Println("length of map: ", len(globalDict.m))
	fmt.Println("Dataset Init Time: ", timer1.Sub(timer0))
	// fmt.Println("Gluster Init Time: ", timer2.Sub(timer1))
	fmt.Println("wordCount Time:    ", timer3.Sub(timer2))
}
