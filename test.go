package main

import(
	"./gluster"
	"fmt"
)

type Blah struct {
	A int
	B int
}

func main() {
	gluster.AddRunner("localhost")

	gluster.ImportFunctionFile("functions/functions.go")

	var ret Blah
	gluster.RunDist("functions.Multiply", &ret, 27, 53)

	fmt.Println("Got back", ret)
}