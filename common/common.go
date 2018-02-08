package common

import(
	"net"
)

/*
* Constants
*/
const DEFAULT_PORT = 1234
const HASH_CMP byte = 11
const SEND_FILE byte = 12
const EXEC_FUNC byte = 13
const ACK byte = 1
const NACK byte = 0

type FuncFile struct {
	CallPrefix string
	Contents []byte
	Checksum []byte
}

type ExecSend struct {
	FuncFile string
	FuncName string
}

func SendByte(conn net.Conn, b byte){
	buf := make([]byte, 1)
	buf[0] = b
	conn.Write(buf)
}

func RecvByte(conn net.Conn) byte{
	buf := make([]byte, 1)
	var _, err = conn.Read(buf)
	if(err != nil){
		return 0
	}
	
	return buf[0]
}

func RecvACK(conn net.Conn) bool{
	if(RecvByte(conn) == 1){
		return true
	}
	return false
}