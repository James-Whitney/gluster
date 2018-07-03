package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var outputMatrix []int

func fillArray(a []int, x int) {
	for i := 0; i < len(a); i++ {
		a[i] = x
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

//MatrixMultiply ...
func MatrixMultiply(inputA []int, inputB []int, width int, id int, idCount int) []int {
	//fmt.Println("Matrix Multiply: ", len(inputA), " x ", len(b))
	fmt.Println("Beginning Matrix Multiplication")
	outputMatrix = make([]int, width*width)
	var start = id * width / idCount
	var end = (id + 1) * width / idCount

	for row := start; row < end; row++ {
		for col := 0; col < width; col++ {
			var Pvalue = 0
			for k := 0; k < width; k++ {
				Pvalue += inputA[row*width+k] * inputB[k*width+col]
			}
			outputMatrix[row*width+col] = Pvalue
		}
	}
	fmt.Println("Matrix Multiplication Complete")

	if(!verifyOutput(outputMatrix, width)){
		fmt.Println("Output array incorrect!!!")
	}

	return outputMatrix
}

func routineMM(inputA []int, inputB []int, width int, row int) {
	for col := 0; col < width; col++ {
		Pvalue := 0
		for k := 0; k < width; k++ {
			Pvalue += inputA[row*width+k] * inputB[k*width+col]
		}
		outputMatrix[row*width+col] = Pvalue
	}
	return
}

func goRoutinesMatrixMultiply(inputA []int, inputB []int, width int) []int {
	outputMatrix = make([]int, width*width)
	for row := 0; row < width; row++ {
		go routineMM(inputA, inputB, width, row)
	}
	//fmt.Println("output: ", outputMatrix)

	if(!verifyOutput(outputMatrix, width)){
		fmt.Println("Output array incorrect!!!")
	}

	return outputMatrix
}

func manyroutineMM(inputA []int, inputB []int, width int, row int, col int) {
	Pvalue := 0
	for k := 0; k < width; k++ {
		Pvalue += inputA[row*width+k] * inputB[k*width+col]
	}
	outputMatrix[row*width+col] = Pvalue
}

func manygoRoutinesMatrixMultiply(inputA []int, inputB []int, width int) []int {
	outputMatrix = make([]int, width*width)
	for row := 0; row < width; row++ {
		for col := 0; col < width; col++ {
			go manyroutineMM(inputA, inputB, width, row, col)
		}
	}

	if(!verifyOutput(outputMatrix, width)){
		fmt.Println("Output array incorrect!!!")
	}

	return outputMatrix
}

func main() {
	ArraySize, _ := strconv.Atoi(os.Args[1])
	//processCount, _ := strconv.Atoi(os.Args[2])

	//create arrays
	matrixA := make([]int, ArraySize*ArraySize)
	fillArray(matrixA, 1)
	matrixB := make([]int, ArraySize*ArraySize)
	fillArray(matrixB, 2)

	fmt.Println("Sequential MM...")
	timer1 := time.Now()
	MatrixMultiply(matrixA, matrixB, ArraySize, 0, 1)
	timer2 := time.Now()

	fmt.Println("GoRoutines MM...")
	timer3 := time.Now()
	goRoutinesMatrixMultiply(matrixA, matrixB, ArraySize)
	timer4 := time.Now()

	fmt.Println("GoRoutines many MM...")
	timer5 := time.Now()
	manygoRoutinesMatrixMultiply(matrixA, matrixB, ArraySize)
	timer6 := time.Now()

	fmt.Println("Sequential time:        ", timer2.Sub(timer1))
	fmt.Println("GoRoutines by row time: ", timer4.Sub(timer3))
	fmt.Println("GoRoutines by ele time: ", timer6.Sub(timer5))

}
