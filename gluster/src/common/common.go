package common

import(
	"net"
	"reflect"
)

/*
* Constants
*/
const DEFAULT_PORT = 1234

//command control bytes 
const EXEC_FUNC byte = 13
const GET_INFO byte = 14

const REQUESTING_FILE byte = 15
const FILE_INCOMING byte = 16
const REQUESTING_ARGS byte = 17
const ARGS_INCOMING byte = 18

const ACK byte = 1
const NACK byte = 0

const GO_FILE = 0
const SO_FILE = 1

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
	FileType uint8
	Content []byte
}

type ExecRequest struct {
	Checksum uint32
	FuncFileName string
	FuncName string
}

type FuncSignature struct {
	Out string
	In []string
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

func encodeTypeHelper(t reflect.Type, structList []string) string{
	//expand struct
	if(t.Kind() == reflect.Struct){
		//check for recursion on structs and throw error
		structList = append(structList, t.Name())
		//encode each field in struct
		var str string = "{"
		var i int
		for i = 0;i < t.NumField() - 1; i++ {
			str += encodeTypeHelper(t.Field(i).Type, structList) + ", "
		}
		if(t.NumField() > 0){
			str += encodeTypeHelper(t.Field(i).Type, structList)
		}
		return str + "}"
	} else if(t.Kind() == reflect.Ptr){
		return "*" + encodeTypeHelper(t.Elem(), structList)
	} else {
		return t.Kind().String()
	}
}

func EncodeType(t reflect.Type) string {
	return encodeTypeHelper(t, nil)
}