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
	var id = gluster.RunDist("functions.Multiply", &ret, 53, 1)

	//wait
	fmt.Println("Waiting")
	for(!gluster.JobDone(id)){}

	fmt.Println("Got back", ret)
}