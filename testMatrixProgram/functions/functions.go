package functions

/*
func MatrixMultiply(a []int, b []int, width int, id int, idCount int) []int {
	//fmt.Println("Matrix Multiply: ", len(a), " x ", len(b))

	c := make([]int, len(a))
	var start = id * len(a) / idCount
	var end = (id + 1) * len(a) / idCount

	for row := start; row < end; row++ {
		for col := 0; col < len(a); col++ {
			var Pvalue = 0
			for k := 0; k < len(a); k++ {
				Pvalue += a[row*len(a)+k] * b[k*len(a)+col]
			}
			c[row*len(a)+col] = Pvalue
		}
	}
	return c
}*/

//MatrixSum ...Gtes the sum of all the values in a matrix
func MatrixSum(a []int, width int, id int, idCount int) int {
	var sum = 0
	var start = id * width / idCount
	var end = (id + 1) * width / idCount
	for row := start; row < end; row++ {
		for col := 0; col < width; col++ {
			sum += a[row*width+col]
		}
	}
	return sum
}
