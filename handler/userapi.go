package handler

import (
	"crypto/md5"
	"encoding/hex"
	"entryTask/constant"
	"entryTask/model"
	"entryTask/rpc"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"

	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//clientToken, err := r.Cookie("token")
	//if err != nil {
	//	log.Printf("get request token err: %v", err)
	//}
	//
	//var ValidateToken func(string) (bool, error)
	//rpcClient := rpc.NewClient(constant.TCP_ADDR)
	//rpcClient.CallRPC("ValidateToken", &ValidateToken)
	//
	//valid, _ := ValidateToken(clientToken.Raw)
	//if valid {
	//	UpdateToken(clientToken.Raw)
	//
	//}

	if r.Method == http.MethodGet {
		loginPage, err := ioutil.ReadFile("./static/login.html")
		if err != nil {
			log.Fatalln("fail to read static/login.html")
		}

		clientToken, err := r.Cookie("token")
		if err != nil {
			log.Printf("get token err: %v", err)
		}

		if clientToken == nil || clientToken.Value == "" {
			log.Printf("token not found")
			w.Write(loginPage)
			return
		}
		rpcClient := GetClient(clientToken.Value)
		var ValidateToken func(string) (bool, error)
		rpcClient.CallRPC("ValidateToken", &ValidateToken)
		valid, _ := ValidateToken(clientToken.Value)
		if valid {
			http.Redirect(w, r, "/profile", http.StatusMovedPermanently)
			return
		}



	} else if r.Method == http.MethodPost {
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		fmt.Println(r.Form)

		// TODO RPC CALL TO authenticate
		rpcClient := rpc.NewClient(constant.TCP_ADDR)

		var LoginAuth func(string, string) (bool, string, error)

		rpcClient.CallRPC("LoginAuth", &LoginAuth)
		ok, msg, err := LoginAuth(username, password)
		if err != nil {
			log.Printf("rpc err: %v", err.Error())
		}
		if !ok {
			respMsg := model.RespMsg{
				Code: "200",
				Msg:  msg,
				Data: nil,
			}

			w.Write(respMsg.JsonBytes())
			return
		}

		var AddTokenByUsername func(string) (bool, string, error)
		rpcClient.CallRPC("AddTokenByUsername", &AddTokenByUsername)
		ok, token, err := AddTokenByUsername(username)

		respMsg := model.RespMsg{
			Code: "200",
			Msg: msg,
			Data: struct {
				Location string
				Token string
			} {
				Location: "http://" + constant.HTTP_ADDR + "/profile",
				Token: token,
			},
		}
		w.Header().Set("Set-Cookie", "token=" + token)
		w.Write(respMsg.JsonBytes())
		return
	}
}


func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	token, _:= r.Cookie("token")
	if r.Method == http.MethodGet{
		var GetUserByToken func(string) (model.User, error)
		client := GetClient(token.Value)
		client.CallRPC("GetUserByToken", &GetUserByToken)
		user, _:= GetUserByToken(token.Value)
		fmt.Println(user.Username, user.Nickname, user.PicUrl)
		profileTemplate := template.Must(template.ParseFiles("./static/profile.html"))
		profileTemplate.Execute(w, user.ToMap())

	} else if r.Method == http.MethodPost {
		client := GetClient(token.Value)
		nickname := r.Form.Get("nickname")

		file, header, err := r.FormFile("profile-picture")
		if err != nil {
			log.Printf("read form-file err: %v", err)
		}
		if file != nil {
			a := md5.Sum([]byte(header.Filename))
			fileNameHash := hex.EncodeToString(a[:])
			localFile, err := os.Create("./picture/" + fileNameHash)
			if err != nil {
				log.Printf("fail to create file, err: %v ", err)
			}
			defer localFile.Close()
			_, err = io.CopyN(localFile, file, header.Size)
			if err != nil {
				log.Printf("fail to copy file, err: %v ", err)
			}
			var UpdateUserPicUrlByToken func(token string, newPicUrl string) (bool, error)
			client.CallRPC("UpdateUserPicUrlByToken", &UpdateUserPicUrlByToken)
			ok, _ := UpdateUserPicUrlByToken(token.Value, "http://" + constant.HTTP_ADDR + "/picture/" + fileNameHash)
			if !ok {
				w.Write([]byte("err"))
				return
			}
		}

		if nickname != "" {
			fmt.Println(nickname)
			var UpdateUserNicknameByToken func(token string, newNickname string) (bool, error)
			client.CallRPC("UpdateUserNicknameByToken", &UpdateUserNicknameByToken)
			ok, _ := UpdateUserNicknameByToken(token.Value, nickname)
			if !ok {
				w.Write([]byte("err"))
				return
			}
		}
		resp := model.RespMsg{
			Code: "200",
			Msg:  "Updated, Redirect",
			Data: struct {
				Location string
			}{
				Location: "http://" + constant.HTTP_ADDR + "/profile",
			},
		}
		w.Write(resp.JsonBytes())
	}
}



func QueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	fmt.Println(r.RemoteAddr)
	client := GetClient(r.RemoteAddr)
	var GetUserByUsername func(string) (model.User, error)
	client.CallRPC("GetUserByUsername", &GetUserByUsername)
	user, _ := GetUserByUsername(username)

	w.Write([]byte(user.Username))

}