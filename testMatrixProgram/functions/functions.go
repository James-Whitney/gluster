package main

//MatrixMultiply ...
func MatrixMultiply(inputA []int, inputB []int, width int, id int, idCount int) []int {
	//fmt.Println("Matrix Multiply: ", len(inputA), " x ", len(b))

	outputMatrix := make([]int, width*width)
	var start = id * len(inputA) / idCount
	var end = (id + 1) * len(inputA) / idCount

	for row := start; row < end; row++ {
		for col := 0; col < len(inputA); col++ {
			var Pvalue = 0
			for k := 0; k < len(inputA); k++ {
				Pvalue += inputA[row*len(inputA)+k] * inputB[k*len(inputA)+col]
			}
			outputMatrix[row*len(inputA)+col] = Pvalue
		}
	}
	return outputMatrix
}

//MatrixSum ... Gets the sum of all the values in inputA matrix
func MatrixSum(inputA []int, width int, id int, idCount int) int {
	var sum = 0
	var start = id * width / idCount
	var end = (id + 1) * width / idCount
	for row := start; row < end; row++ {
		for col := 0; col < width; col++ {
			sum += inputA[row*width+col]
		}
	}
	return sum
}