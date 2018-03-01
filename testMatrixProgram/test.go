package main

import (
	"fmt"
	"reflect"
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

func testMatrixSum() {
	const maxArraySize int = 8
	const processCount int = 1
	const expected int = 192

	fmt.Println("Getting the Sum of a Matrix...")

	inputArray := make([]int, maxArraySize*maxArraySize)
	fillArray(inputArray, 3)
	printMatrix(inputArray, maxArraySize)

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

func testMatrixMultiplication() {
	const maxArraySize int = 8
	const processCount int = 1

	fmt.Println("Multiplying Two Matrices...")

	fmt.Println("Matrix A:")
	inputA := make([]int, maxArraySize*maxArraySize)
	fillArray(inputA, 3)
	printMatrix(inputA, maxArraySize)

	fmt.Println("Matrix B:")
	inputB := make([]int, maxArraySize*maxArraySize)
	fillArray(inputB, 4)
	printMatrix(inputB, maxArraySize)

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
	fmt.Println("Result Matrix:")
	printMatrix(output, maxArraySize)
}

func main() {
	timer1 := time.Now()

	gluster.AddRunner("localhost")
	gluster.ImportFunctionFile("functions/functions.go")

	timer2 := time.Now()

	testMatrixSum()

	timer3 := time.Now()

	testMatrixMultiplication()

	timerEnd := time.Now()

	fmt.Println("Gluster Init Time: 	  ", timer2.Sub(timer1))
	fmt.Println("testMatrixSum Time:   ", timer3.Sub(timer2))
	fmt.Println("testMatrixMulti Time: ", timerEnd.Sub(timer3))

}
