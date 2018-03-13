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

func verifyOutput(output []int, width int) bool {
	for i := range output {
		if(i != 2*width){
			return false
		}
	}
	return true;
}

func testMatrixSum(maxArraySize int, processCount int) {
	const expected int = 192

	fmt.Println("Getting the Sum of a Matrix...")

	inputArray := make([]int, maxArraySize*maxArraySize)
	fillArray(inputArray, 1)
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
	fillArray(inputA, 1)

	inputB := make([]int, maxArraySize*maxArraySize)
	fillArray(inputB, 2)

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

	if(!verifyOutput(output, maxArraySize)){
		fmt.Println("Output array incorrect!!!")
	}
}

func testRoutinesMultiplication(maxArraySize int, processCount int) {
	inputA := make([]int, maxArraySize*maxArraySize)
	fillArray(inputA, 1)
	inputB := make([]int, maxArraySize*maxArraySize)
	fillArray(inputB, 2)

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
	fillArray(inputA, 1)
	inputB := make([]int, maxArraySize*maxArraySize)
	fillArray(inputB, 2)

	output := make([]int, maxArraySize*maxArraySize)

	var runnerList []int

	for i := 0; i < processCount; i++ {
		fmt.Println("Launching Runner: ", i)
		runnerList = append(runnerList, gluster.RunDist("functions.ManygoRoutinesMatrixMultiply", reflect.TypeOf(output), inputA, inputB, maxArraySize, i, processCount))
	}

	for _, runner := range runnerList {
		for !(gluster.JobDone(runner)) {
		}
		var partialOutput = gluster.GetReturn(runner).([]int)
		mergeArray(output, partialOutput)
	}

	if(!verifyOutput(output, maxArraySize)){
		fmt.Println("Output array incorrect!!!")
	}
}

func addRunners() {

	gluster.AddRunner("127x01.csc.calpoly.edu")
	gluster.AddRunner("127x02.csc.calpoly.edu")
	//gluster.AddRunner("127x03.csc.calpoly.edu")
	gluster.AddRunner("127x04.csc.calpoly.edu")
	gluster.AddRunner("127x05.csc.calpoly.edu")
	gluster.AddRunner("127x06.csc.calpoly.edu")
	gluster.AddRunner("127x07.csc.calpoly.edu")
	gluster.AddRunner("127x08.csc.calpoly.edu")
	gluster.AddRunner("127x09.csc.calpoly.edu")
	gluster.AddRunner("127x10.csc.calpoly.edu")
	gluster.AddRunner("127x11.csc.calpoly.edu")
	gluster.AddRunner("127x12.csc.calpoly.edu")
	gluster.AddRunner("127x13.csc.calpoly.edu")
	gluster.AddRunner("127x14.csc.calpoly.edu")
	gluster.AddRunner("127x15.csc.calpoly.edu")
	gluster.AddRunner("127x16.csc.calpoly.edu")
	gluster.AddRunner("127x17.csc.calpoly.edu")
	gluster.AddRunner("127x18.csc.calpoly.edu")
	gluster.AddRunner("127x19.csc.calpoly.edu")
	gluster.AddRunner("127x20.csc.calpoly.edu")
	gluster.AddRunner("127x21.csc.calpoly.edu")
	gluster.AddRunner("127x22.csc.calpoly.edu")
	gluster.AddRunner("127x23.csc.calpoly.edu")
	//gluster.AddRunner("127x24.csc.calpoly.edu")
	gluster.AddRunner("127x25.csc.calpoly.edu")
	gluster.AddRunner("127x26.csc.calpoly.edu")
	gluster.AddRunner("127x27.csc.calpoly.edu")
	gluster.AddRunner("127x28.csc.calpoly.edu")
	gluster.AddRunner("127x29.csc.calpoly.edu")
	gluster.AddRunner("127x30.csc.calpoly.edu")
	gluster.AddRunner("127x31.csc.calpoly.edu")
	gluster.AddRunner("127x32.csc.calpoly.edu")
	gluster.AddRunner("127x33.csc.calpoly.edu")
	//gluster.AddRunner("127x34.csc.calpoly.edu")
	gluster.AddRunner("127x35.csc.calpoly.edu")
	gluster.AddRunner("127x36.csc.calpoly.edu")
	gluster.AddRunner("127x37.csc.calpoly.edu")
}

func main() {
	ArraySize, _ := strconv.Atoi(os.Args[1])
	processCount, _ := strconv.Atoi(os.Args[2])

	timer1 := time.Now()
	addRunners()
	gluster.AddRunner("127x35.csc.calpoly.edu")
	gluster.AddRunner("127x34.csc.calpoly.edu")
	gluster.AddRunner("127x33.csc.calpoly.edu")
	gluster.AddRunner("127x32.csc.calpoly.edu")
	gluster.AddRunner("127x31.csc.calpoly.edu")
	gluster.AddRunner("127x30.csc.calpoly.edu")
	gluster.AddRunner("127x29.csc.calpoly.edu")
	gluster.AddRunner("127x28.csc.calpoly.edu")
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
