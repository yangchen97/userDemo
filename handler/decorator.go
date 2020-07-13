package handler

import (
	"entryTask/constant"
	"entryTask/rpc"
	"log"
	"net/http"
)

func Authorized(h func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		clientToken, err := r.Cookie("token")
		if err != nil {
			log.Printf("get token err: %v", err)
		}

		if clientToken == nil || clientToken.Value == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		rpcClient := rpc.NewClient(constant.TCP_ADDR)
		var ValidateToken func(string) (bool, error)
		var UpdateToken func (string) (bool, error)
		rpcClient.CallRPC("ValidateToken", &ValidateToken)
		rpcClient.CallRPC("UpdateToken", &UpdateToken)

		valid, _ := ValidateToken(clientToken.Value)
		if valid {
			h(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	})
}