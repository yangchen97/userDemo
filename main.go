package main

import (
	"encoding/gob"
	"entryTask/constant"
	"entryTask/handler"
	"entryTask/model"
	"entryTask/rpc"
	"entryTask/service"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	var start string

	flag.StringVar(&start, "type",start, "tcp or http")
	flag.Parse()
	fmt.Println(start)

	if start == "tcp" {
		startTcp()
	} else if start == "http" {
		startHttp()
	}


}

func startHttp() {
	gob.Register(model.User{})
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/profile", handler.Authorized(handler.ProfileHandler))
	http.Handle("/picture/", http.StripPrefix("/picture/", http.FileServer(http.Dir("picture"))))
	http.ListenAndServe(constant.HTTP_ADDR, nil)
}

func startTcp() {
	gob.Register(model.User{})
	server := rpc.NewServer(constant.TCP_ADDR)
	//server.Register("QueryUserByName", handler.QueryUserByName)
	//server.Register("Login", handler.Login)
	//server.Register("ValidateToken", handler.ValidateToken)
	//server.Register("ModifyNickname", handler.ModifyNickname)
	server.Register("LoginAuth", service.LoginAuth)
	server.Register("AddTokenByUsername", service.AddTokenByUsername)
	server.Register("GetUserByToken", service.GetUserByToken)
	server.Register("ValidateToken",service.ValidateToken)
	server.Register("UpdateToken",service.UpdateToken)
	server.Register("UpdateUserNicknameByToken", service.UpdateUserNicknameByToken)
	server.Register("UpdateUserPicUrlByToken", service.UpdateUserPicUrlByToken)

	fmt.Println("RPC server running...")
	server.Run()
}
