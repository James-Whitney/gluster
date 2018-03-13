package main

import (
	"fmt"
	"sync"
)

func printMatrix(x []int, size int) {
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			fmt.Print(x[i*size+j], " ")
		}
		fmt.Println()
	}
}

func ManygoRoutinesMatrixMultiply(inputA []int, inputB []int, width int, machineID int, machineCount int) []int {
	var wg sync.WaitGroup
	outputMatrix := make([]int, width*width)

	var start = machineID * width / machineCount
	var end = (machineID + 1) * width / machineCount

	for row := start; row < end; row++ {
		for col := 0; col < width; col++ {
			wg.Add(1)
			go func(rowMM int, colMM int) {
				Pvalue := 0
				for k := 0; k < width; k++ {
					Pvalue += inputA[rowMM*width+k] * inputB[k*width+colMM]
				}
				outputMatrix[rowMM*width+colMM] = Pvalue
				defer wg.Done()
			}(row, col)
		}
	}
	wg.Wait()
	return outputMatrix
}

func RoutinesMatrixMultiply(inputA []int, inputB []int, width int, machineID int, machineCount int) []int {
	var wg sync.WaitGroup
	outputMatrix := make([]int, width*width)

	var start = machineID * width / machineCount
	var end = (machineID + 1) * width / machineCount

	for row := start; row < end; row++ {
		wg.Add(1)
		go func(rowMM int) {
			for col := 0; col < width; col++ {
				Pvalue := 0
				for k := 0; k < width; k++ {
					Pvalue += inputA[row*width+k] * inputB[k*width+col]
				}
				outputMatrix[row*width+col] = Pvalue
			}
			defer wg.Done()
		}(row)
	}
	wg.Wait()
	return outputMatrix
}

//MatrixMultiply ...
func MatrixMultiply(inputA []int, inputB []int, width int, id int, idCount int) []int {
	//fmt.Println("Matrix Multiply: ", len(inputA), " x ", len(b))
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
