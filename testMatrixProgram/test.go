package main

import (
	"fmt"
	"reflect"

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
	const processCount int = 10
	const expected int = 192

	fmt.Println("Getting the Sum of a Matrix...")

	inputArray := make([]int, maxArraySize*maxArraySize)
	fillArray(inputArray, 3)
	printMatrix(inputArray, maxArraySize)

	var sum int

	//LAUNCH SLAVES
	for i := 0; i < processCount; i++ {
		fmt.Println("Launching Runner: ", i)
		gluster.RunDist("functions.MatrixSum", reflect.TypeOf(sum), inputArray, maxArraySize, i, processCount)
		//c = MatrixMultiply(a, b, maxArraySize, i, processCount)
	}

	//WAIT FOR RETURNS FROM SLAVES
	for i := 0; i < processCount; i++ {
		for !(gluster.JobDone(i)) {
		}
		var particalSum = gluster.GetReturn(i).(int)
		fmt.Println("Waited for: ", i, " Partial Sum: ", particalSum)
		sum += particalSum
	}
	fmt.Print("Expected: ", expected, " Actual: ", sum)
	if sum == expected {
		fmt.Println(" SUCCESS")
	} else {
		fmt.Println(" FAILURE")
	}
}

func main() {
	gluster.AddRunner("localhost")
	gluster.ImportFunctionFile("functions/functions.go")

	testMatrixSum()

	//fillArray(b, 2)
	//returnArray := make([][]int, processCount)

	//result := make([]int, maxArraySize*maxArraySize)

	//wait
	/*
		fmt.Println("Waiting")
		for !gluster.JobDone(id) {
		}

		fmt.Println("Got back", ret)
	*/
}
