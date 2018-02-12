package main

import "fmt"

type Blah struct {
	A int
	B int
}

func Multiply(a, b int) Blah{
	//*reply = args.A * args.B
	fmt.Println("test", a, b)

	return Blah{a+128, a * b +1}
}