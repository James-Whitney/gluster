package main

import (
	"fmt"

	"./gluster"
)

type Blah struct {
	A int
	B int
}

const maxArraySize int = 8

func fillArray(a []int, x int) {
	for i := 0; i < len(a); i++ {
		a[i] = x
	}
}

func mergeArray(c []int, d []int) {
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
	b := make([]int, maxArraySize*maxArraySize)
	fillArray(a, 3)
	fillArray(b, 2)
	c := make([]int, maxArraySize*maxArraySize)
	result := make([]int, maxArraySize*maxArraySize)

	for i := 0; i < 10; i++ {
		gluster.RunDist("functions.MatrixMultiply", &c, a, b, i, 10)
		mergeArray(c, result)
	}

	printMatrix(result, maxArraySize)

	var ret Blah
	var id = gluster.RunDist("functions.Multiply", &ret, 53, 1)

	//wait
	fmt.Println("Waiting")
	for !gluster.JobDone(id) {
	}

	fmt.Println("Got back", ret)
}
