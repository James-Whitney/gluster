package main

import (
	"fmt"
)

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

func RoutinesMatrixMultiply(inputA []int, inputB []int, width int, id int, idCount int) []int {
	fmt.Println("Beginning Routines Matrix Multiplication")
	outputMatrix := make([]int, width*width)
	var start = id * width / idCount
	var end = (id + 1) * width / idCount

	var channels []chan int

	for row := start; row < end; row++ {
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
	fmt.Println("Matrix Routines Multiplication Complete")

	return outputMatrix
}

//MatrixMultiply ...
func MatrixMultiply(inputA []int, inputB []int, width int, id int, idCount int) []int {
	//fmt.Println("Matrix Multiply: ", len(inputA), " x ", len(b))
	fmt.Println("Beginning Matrix Multiplication")
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
