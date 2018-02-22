package main

import (
	"fmt"

	"../gluster/src/master"
)

const maxArraySize int = 8
const processCount int = 10

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

func main() {
	gluster.AddRunner("localhost")
	gluster.ImportFunctionFile("functions/functions.go")

	a := make([]int, maxArraySize*maxArraySize)
	//b := make([]int, maxArraySize*maxArraySize)
	fillArray(a, 3)
	printMatrix(a, maxArraySize)
	var glusterResults [processCount]int
	var result int = 0
	//fillArray(b, 2)
	//returnArray := make([][]int, processCount)

	//result := make([]int, maxArraySize*maxArraySize)

	for i := 0; i < processCount; i++ {
		//c := make([]int, maxArraySize*maxArraySize)
		fmt.Println("Launching Runner: ", i)
		gluster.RunDist("functions.MatrixSum", &glusterResults[i], a, maxArraySize, i, processCount)
		//c = MatrixMultiply(a, b, maxArraySize, i, processCount)
	}

	//wait

	for i := 0; i < processCount; i++ {
		for !(gluster.JobDone(i)) {
		}
		fmt.Println("Waited for: ", i, " result: ", glusterResults[i])
		result += glusterResults[i]
	}

	fmt.Println("Result of MatrixSum: ", result)

	//wait
	/*
		fmt.Println("Waiting")
		for !gluster.JobDone(id) {
		}

		fmt.Println("Got back", ret)
	*/
}
