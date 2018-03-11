package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func fillArray(a []int, x int) {
	for i := 0; i < len(a); i++ {
		a[i] = x
	}
}

//MatrixMultiply ...
func MatrixMultiply(inputA []int, inputB []int, width int, id int, idCount int) []int {
	//fmt.Println("Matrix Multiply: ", len(inputA), " x ", len(b))
	fmt.Println("Beginning Matrix Multiplication")
	//fmt.Println("InputA: ", inputA)

	//fmt.Println()

	//fmt.Println("InputB: ", inputB)
	outputMatrix := make([]int, width*width)
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
	return outputMatrix
}

func routineMM(inputA []int, inputB []int, width int, row int, channel chan<- int) {
	for col := 0; col < width; col++ {
		Pvalue := 0
		for k := 0; k < width; k++ {
			Pvalue += inputA[row*width+k] * inputB[k*width+col]
		}
		channel <- Pvalue
	}
	close(channel)
	return
}

func goRoutinesMatrixMultiply(inputA []int, inputB []int, width int) []int {
	var channels []chan int

	outputMatrix := make([]int, width*width)
	for row := 0; row < width; row++ {
		ch := make(chan int)
		channels = append(channels, ch)
		go routineMM(inputA, inputB, width, row, ch)
		//outputMatrix[row*width+col] = <-channels[row*width+col]
	}
	for row := 0; row < width; row++ {
		for col := 0; col < width; col++ {
			outputMatrix[row*width+col] = <-channels[row]
		}
	}
	//fmt.Println("output: ", outputMatrix)
	return outputMatrix
}

func manyroutineMM(inputA []int, inputB []int, width int, row int, col int, channel chan<- int) {
	Pvalue := 0
	for k := 0; k < width; k++ {
		Pvalue += inputA[row*width+k] * inputB[k*width+col]
	}
	channel <- Pvalue
	close(channel)
}

func manygoRoutinesMatrixMultiply(inputA []int, inputB []int, width int) []int {
	var channels []chan int

	outputMatrix := make([]int, width*width)
	for row := 0; row < width; row++ {
		for col := 0; col < width; col++ {
			ch := make(chan int)
			channels = append(channels, ch)
			go manyroutineMM(inputA, inputB, width, row, col, ch)
			//outputMatrix[row*width+col] = <-channels[row*width+col]
		}
	}
	for row := 0; row < width; row++ {
		for col := 0; col < width; col++ {
			outputMatrix[row*width+col] = <-channels[row*width+col]
		}
	}
	//fmt.Println("output: ", outputMatrix)
	return outputMatrix
}

func main() {
	ArraySize, _ := strconv.Atoi(os.Args[1])
	//processCount, _ := strconv.Atoi(os.Args[2])

	//create arrays
	matrixA := make([]int, ArraySize*ArraySize)
	fillArray(matrixA, 3)
	matrixB := make([]int, ArraySize*ArraySize)
	fillArray(matrixB, 4)

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

	fmt.Println("Sequential time:       ", timer2.Sub(timer1))
	fmt.Println("GoRoutines time:       ", timer4.Sub(timer3))
	fmt.Println("GoRoutines many time:  ", timer6.Sub(timer5))

}
