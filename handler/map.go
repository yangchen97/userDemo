package handler

import (
	"entryTask/constant"
	"entryTask/rpc"
	"fmt"
	"sync"
)

var clientMap sync.Map

func init() {

}


func GetClient(remoteAdrr string) *rpc.Client {
	client, ok := clientMap.Load(remoteAdrr)
	if !ok {
		fmt.Println("not found")
		newClient := rpc.NewClient(constant.TCP_ADDR)
		clientMap.Store(remoteAdrr, newClient)
		return newClient
	}
	fmt.Println("found")
	return client.(*rpc.Client)
}