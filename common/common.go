package common

/*
* Constants
*/
const DEFAULT_PORT = 1234
const HASH_CMP byte = 11
const SEND_FILE byte = 12
const EXEC_FUNC byte = 13

type FuncFile struct {
	CallPrefix string
	Contents []byte
	Checksum []byte
}

type ExecSend struct {
	FuncFile string
	FuncName string
	Args interface{}
}

type ExecResp struct {
	Reply interface{}
}