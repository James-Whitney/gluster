package main

import "fmt"

type Blah struct {
	A int
	B int
}

func Multiply(a, b int) Blah {
	//*reply = args.A * args.B
	fmt.Println("test", a, b)

	return Blah{a + 128, a*b + 1}
}

func MatrixMultiply(a []int, b []int, id int, idCount int) []int {
	fmt.Println("Matrix Multiply: ", len(a), " x ", len(b))

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
}
