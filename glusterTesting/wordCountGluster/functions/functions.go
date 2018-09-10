package main

import (
	"sync"
	"strings"
)

func WordCount(words []string, processCount int) map[string]int {
	// fmt.Println("Begining WordCount")
	// perform wordCount
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
			// fmt.Println("Runninng Processes", processID)
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
	// fmt.Println("Ending WordCount")
	// fmt.Println("Length of map: ", len(globalDict.m))
	return globalDict.m
}