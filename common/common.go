package common

import(
	"net"
)

/*
* Constants
*/
const DEFAULT_PORT = 1234

//command control bytes
const HASH_CMP byte = 11
const SEND_FILE byte = 12
const EXEC_FUNC byte = 13
const GET_INFO byte = 14

const ACK byte = 1
const NACK byte = 0

type RunnerInfo struct {
	Cores int
	Arch string
	OS string
}

type RunnerStatus struct {
	LoadedPlugins []string
	SystemLoad float32
	RunningJobs int
}

type FuncFile struct {
	CallPrefix string
	Checksum uint32
}

type FuncFileContent struct {
	File FuncFile
	Content []byte
}

type ExecSend struct {
	FuncFileName string
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