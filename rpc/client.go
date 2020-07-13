package rpc

import (
	"entryTask/constant"
	. "entryTask/transport"
	"fmt"
	"log"
	"net"
	"reflect"
)

type Client struct {
	addr string
	conn net.Conn
}

func NewClient(addr string) *Client {
	conn, err := net.Dial("tcp", constant.TCP_ADDR)
	if err != nil {
		log.Printf("new rpc client err: %v", err)
	}
	return &Client{
		addr: addr,
		conn: conn,
	}
}

func (client *Client) CallRPC(funcName string, funcPtr interface{}) {
	container := reflect.ValueOf(funcPtr).Elem() // funcPtr指向一个函数，elem()解引用
	f := func(request []reflect.Value) []reflect.Value {
		transport := NewTransport(client.conn)

		// 提取参数包装成RPC数据
		inArgs := make([]interface{}, 0, len(request))
		for _, arg := range request {
			inArgs = append(inArgs, arg.Interface())
		}
		requestRPCdata := RPCdata{
			FuncName: funcName,
			Args:     inArgs,
		}

		// 参数编码
		requestTLV, err := Encode(requestRPCdata)
		if err != nil {
			panic(err)
		}

		err = transport.Send(requestTLV)
		if err != nil {
			panic(err)
		}

		responseByte, err := transport.Read()
		if err != nil {
			panic(err)
		}

		response, err := Decode(responseByte)
		if err != nil {
			panic(err)
		}

		if len(response.Args) == 0 {
			response.Args = make([]interface{}, container.Type().NumOut()) //numOut:反射函数的返回结果的数量
		}
		numOut := container.Type().NumOut()
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut-1 {
				if response.Args[i] == nil {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(response.Args[i])
				}
			} else {
				if response.Error == "" {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(fmt.Errorf(response.Error))
				}
			}
		}

		return outArgs

	}
	container.Set(reflect.MakeFunc(container.Type(), f))

}
