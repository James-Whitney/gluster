package main

import "fmt"

type Blah struct {
	A int
	B int
}

func Multiply(args interface{}, reply interface{}) {
	//*reply = args.A * args.B
	fmt.Println("test")
}