package main

import "fmt"

type Blah struct {
	A int
	B int
}

func Multiply(a, b int, q Blah) {
	//*reply = args.A * args.B
	fmt.Println("test", a, b, q)
}