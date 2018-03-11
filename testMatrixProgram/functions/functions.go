package main

import (
	"fmt"
)

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
