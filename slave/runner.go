package main

import (
	"net"
	"fmt"
	"net/rpc"
	"../functions"
)

func main() {
	handle_work()
}



/*
func wait_for_work()
{
	//listen for connection from master
	listen , err := net.Listen("tcp", "localhost:8081");
	if err != nil {
		fmt.Println("Error setting up tcp connection: ", err);
		return;
	}
	defer listen.Close();

	//loop waiting for work connections
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err);
			return;
		}

		//thread work to be done
		go handle_work(conn);
	}
}*/


func handle_work() {
	//defer conn.Close();

	distFuncs := new(functions.DistFunc)
	rpc.Register(distFuncs)

	listen, err := net.Listen("tcp", ":1234") //TODO choose port dynamically
	if err != nil {
		fmt.Println("Error setting up tcp listener: ", err)
	}

	rpc.Accept(listen)
}