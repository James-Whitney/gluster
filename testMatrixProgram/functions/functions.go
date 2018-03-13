package main

import (
	"fmt"
)

var outputMatrix []int

func printMatrix(x []int, size int) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Print(x[i*size+j], " ")
		}
		fmt.Println()
	}
}

func manyroutineMM(inputA []int, inputB []int, width int, row int, col int) {
	Pvalue := 0
	for k := 0; k < width; k++ {
		Pvalue += inputA[row*width+k] * inputB[k*width+col]
	}
	outputMatrix[row*width+col] = Pvalue
}

func ManygoRoutinesMatrixMultiply(inputA []int, inputB []int, width int, machineID int, machineCount int) []int {
	outputMatrix = make([]int, width*width)
	var start = machineID * width / machineCount
	var end = (machineID + 1) * width / machineCount
	for row := start; row < end; row++ {
		for col := 0; col < width; col++ {
			go manyroutineMM(inputA, inputB, width, row, col)
		}
	}
	printMatrix(outputMatrix, width)
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

func RoutinesMatrixMultiply(inputA []int, inputB []int, width int, machineID int, machineCount int) []int {
	fmt.Println("Beginning Routines Matrix Multiplication")

	outputMatrix = make([]int, width*width)
	var start = machineID * width / machineCount
	var end = (machineID + 1) * width / machineCount

	for row := start; row < end; row++ {
		go routineMM(inputA, inputB, width, row)
	}
	fmt.Println("Matrix Routines Multiplication Complete")
	printMatrix(outputMatrix, width)
	return outputMatrix
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
	printMatrix(outputMatrix, width)
	return outputMatrix
}

//MatrixSum ... Gets the sum of all the values in inputA matrix
func MatrixSum(inputA []int, width int, id int, idCount int) int {
	fmt.Println("Beginning Matrix Sum")
	var sum = 0
	var start = id * width / idCount
	var end = (id + 1) * width / idCount
	for row := start; row < end; row++ {
		for col := 0; col < width; col++ {
			sum += inputA[row*width+col]
		}
	}
	fmt.Println("Matrix Sum Complete")
	return sum
}
