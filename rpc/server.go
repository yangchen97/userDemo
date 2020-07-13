package rpc

import (
	. "entryTask/transport"
	"fmt"
	"log"
	"net"
	"reflect"
)

type Server struct {
	addr      string
	functions map[string]reflect.Value
}

func NewServer(addr string) *Server {
	functions := make(map[string]reflect.Value)
	return &Server{addr: addr, functions: functions}
}

func (server *Server) Register(funcName string, function interface{}) {
	_, ok := server.functions[funcName]
	if ok {
		return
	}
	server.functions[funcName] = reflect.ValueOf(function)
}

func (server *Server) Execute(request RPCdata) RPCdata {
	function, ok := server.functions[request.FuncName]
	if !ok {
		e := fmt.Sprintf("function %s is not registered", request.FuncName)
		log.Println(e)
		return RPCdata{
			FuncName: request.FuncName,
			Args:     nil,
			Error:    e,
		}
	}

	// 提取参数并运行函数
	args := make([]reflect.Value, len(request.Args))
	for i, arg := range request.Args {
		args[i] = reflect.ValueOf(arg)
	}
	result := function.Call(args)

	// function返回的最后一个参数为error，先去掉
	resultArgs := make([]interface{}, len(result)-1)
	for i, arg := range result {
		if i == len(result)-1 {
			continue
		}
		resultArgs[i] = arg.Interface()
	}

	// 把函数返回的错误信息提出来
	var err string
	if e, ok := result[len(result)-1].Interface().(error); ok { //.(error)类型强制转换
		// convert the error into error string value
		err = e.Error()
	}
	return RPCdata{
		FuncName: request.FuncName,
		Args:     resultArgs,
		Error:    err,
	}

}

func (server *Server) Run() {
	localAddress, _ := net.ResolveTCPAddr("tcp4", server.addr)
	tcpListener, err := net.ListenTCP("tcp", localAddress)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = tcpListener.Close()
	}()
	for {
		fmt.Println("waiting for connection...")
		var conn, err = tcpListener.AcceptTCP() //接受连接
		if err != nil {
			fmt.Println("fail to connect：", err)
			return
		}
		fmt.Println("connect successfully")
		go func() {
			transport := NewTransport(conn)
			for {
				requestByte, err := transport.Read()

				if err != nil {
					log.Println(err)
					break
				}
				requestRPCdata, err := Decode(requestByte)
				if err != nil {
					log.Println(err)
					break
				}

				fmt.Println("server executing")
				response := server.Execute(requestRPCdata)
				responseByte, err := Encode(response)
				if err != nil {
					panic(err)
				}
				transport.Send(responseByte)
			}
		}()
	}
}
