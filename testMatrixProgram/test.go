package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"../gluster/src/master"
)

func fillArray(a []int, x int) {
	for i := 0; i < len(a); i++ {
		a[i] = x
	}
}

func mergeArray(d []int, c []int) {
	for i := 0; i < len(c); i++ {
		d[i] += c[i]
	}
}

func printMatrix(x []int, size int) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Print(x[i*size+j], " ")
		}
		fmt.Println()
	}
}

func testMatrixSum(maxArraySize int, processCount int) {
	const expected int = 192

	fmt.Println("Getting the Sum of a Matrix...")

	inputArray := make([]int, maxArraySize*maxArraySize)
	fillArray(inputArray, 3)
	//printMatrix(inputArray, maxArraySize)

	var runnerList []int

	var sum int

	//LAUNCH SLAVES
	for i := 0; i < processCount; i++ {
		fmt.Println("Launching Runner: ", i)
		runnerList = append(runnerList, gluster.RunDist("functions.MatrixSum", reflect.TypeOf(sum), inputArray, maxArraySize, i, processCount))
	}

	for _, runner := range runnerList {
		for !(gluster.JobDone(runner)) {
		}
		var particalSum = gluster.GetReturn(runner).(int)
		fmt.Println("Waited for: ", runner, " Partial Sum: ", particalSum)
		sum += particalSum
	}

	fmt.Print("Expected: ", expected, " Actual: ", sum)
	if sum == expected {
		fmt.Println(" SUCCESS")
	} else {
		fmt.Println(" FAILURE")
	}
}

func testMatrixMultiplication(maxArraySize int, processCount int) {

	fmt.Println("Multiplying Two Matrices...")

	inputA := make([]int, maxArraySize*maxArraySize)
	fillArray(inputA, 3)

	inputB := make([]int, maxArraySize*maxArraySize)
	fillArray(inputB, 4)

	output := make([]int, maxArraySize*maxArraySize)

	var runnerList []int

	for i := 0; i < processCount; i++ {
		fmt.Println("Launching Runner: ", i)
		runnerList = append(runnerList, gluster.RunDist("functions.MatrixMultiply", reflect.TypeOf(output), inputA, inputB, maxArraySize, i, processCount))
	}

	for _, runner := range runnerList {
		for !(gluster.JobDone(runner)) {
		}
		var partialOutput = gluster.GetReturn(runner).([]int)
		mergeArray(output, partialOutput)
	}
	//fmt.Println("Result Matrix:")
	//printMatrix(output, maxArraySize)
}

func testRoutinesMultiplication(maxArraySize int, processCount int) {
	inputA := make([]int, maxArraySize*maxArraySize)
	fillArray(inputA, 3)
	inputB := make([]int, maxArraySize*maxArraySize)
	fillArray(inputB, 4)

	output := make([]int, maxArraySize*maxArraySize)

	var runnerList []int

	for i := 0; i < processCount; i++ {
		fmt.Println("Launching Runner: ", i)
		runnerList = append(runnerList, gluster.RunDist("functions.RoutinesMatrixMultiply", reflect.TypeOf(output), inputA, inputB, maxArraySize, i, processCount))
	}

	for _, runner := range runnerList {
		for !(gluster.JobDone(runner)) {
		}
		var partialOutput = gluster.GetReturn(runner).([]int)
		mergeArray(output, partialOutput)
	}
}

func testManyRoutinesMultiplication(maxArraySize int, processCount int) {
	inputA := make([]int, maxArraySize*maxArraySize)
	fillArray(inputA, 3)
	inputB := make([]int, maxArraySize*maxArraySize)
	fillArray(inputB, 4)

	output := make([]int, maxArraySize*maxArraySize)

	var runnerList []int

	for i := 0; i < processCount; i++ {
		fmt.Println("Launching Runner: ", i)
		runnerList = append(runnerList, gluster.RunDist("functions.manygoRoutinesMatrixMultiply", reflect.TypeOf(output), inputA, inputB, maxArraySize, i, processCount))
	}

	for _, runner := range runnerList {
		for !(gluster.JobDone(runner)) {
		}
		var partialOutput = gluster.GetReturn(runner).([]int)
		mergeArray(output, partialOutput)
	}
}

func main() {
	ArraySize, _ := strconv.Atoi(os.Args[1])
	processCount, _ := strconv.Atoi(os.Args[2])

	timer1 := time.Now()
	gluster.AddRunner("127x35.csc.calpoly.edu")
	gluster.AddRunner("127x34.csc.calpoly.edu")
	gluster.ImportFunctionFile("functions/functions.go")
	timer2 := time.Now()

	timer3 := time.Now()
	testMatrixMultiplication(ArraySize, processCount)
	timer4 := time.Now()

	timer5 := time.Now()
	testRoutinesMultiplication(ArraySize, processCount)
	timer6 := time.Now()

	timer7 := time.Now()
	testManyRoutinesMultiplication(ArraySize, processCount)
	timer8 := time.Now()

	fmt.Println("Gluster Init Time: 	      ", timer2.Sub(timer1))
	fmt.Println("testMatrixMulti Time:     ", timer4.Sub(timer3))
	fmt.Println("testMatrixMulti+goR Time: ", timer6.Sub(timer5))
	fmt.Println("testMatrixMulti+ManygoR Time: ", timer8.Sub(timer7))

}
