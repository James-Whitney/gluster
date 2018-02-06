package main

import(
	"./gluster"
)

func main() {
	gluster.AddRunner("localhost")

	gluster.ImportFunctionFile("functions/functions.go")

	gluster.RunDist("functions.Multiply", nil, nil)

}