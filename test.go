package main

import(
	"fmt"
	"./gluster"
	. "./functions"
)

func main() {
	gluster.AddRunner("localhost")

	var reply int
	args := Blah{}
	args.A = 1
	args.B = 2
	fmt.Println("hi")
	gluster.RunDist("DistFunc.Multiply", args, &reply)

	fmt.Println(reply)
}