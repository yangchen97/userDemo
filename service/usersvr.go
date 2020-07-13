package service

import (
	"crypto/md5"
	"encoding/hex"
	"entryTask/dao"
	"entryTask/model"
	"log"
	"strconv"
	"time"
)

func LoginAuth(username string, password string) (bool, string, error) {
	user, err := dao.GetUserByUsername(username)

	if err != nil {
		return false, err.Error(), err

	}

	if user == nil {
		return false, "username not found", nil
	}

	if password == user.Password {
		return true, "", nil
	} else {
		return false, "wrong pasword", nil
	}

}

func ValidateToken(token string) (bool, error) {
	value := dao.Get("user" + token)
	if value == "1" {
		return true, nil
	}
	return false, nil
}



func AddTokenByUsername(username string) (bool, string, error) {
	token := genToken(username)
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		log.Printf("database get user err: %v", err)
	}
	dao.AddToken(token, user)
	return true, token, nil
}


func UpdateToken(token string) (bool, error) {
	dao.UpdateToken(token)
	return true, nil
}
func GetUserByToken(token string) (model.User, error) {
	return dao.GetUserByToken(token), nil
}

func UpdateUserNicknameByToken(token string, newNickname string) (bool, error) {
	username := dao.GetUserByToken(token).Username

	ok := dao.UpdateUserNicknameByUsername(username, newNickname)
	if !ok {
		log.Printf("fail to update database")
		return false, nil
	}

	user, err := dao.GetUserByUsername(username)
	if err != nil {
		log.Printf("fail to get database user")
		return false, nil
	}
	dao.UpdateUser(token, user)

	return true, nil

}


func UpdateUserPicUrlByToken(token string, newPicUrl string) (bool, error) {
	username := dao.GetUserByToken(token).Username

	ok := dao.UpdateUserPicUrlByUsername(username, newPicUrl)
	if !ok {
		log.Printf("fail to update database")
		return false, nil
	}
	user, err := dao.GetUserByUsername(username)
	if err != nil {
		log.Printf("fail to get database user")
		return false, nil
	}
	dao.UpdateUser(token, user)
	return true, nil

}

func genToken(username string) string{
	ts := time.Now().Unix()
	data := username + strconv.FormatInt(ts, 10) + "1234"
	h := md5.Sum([]byte(data))
	return hex.EncodeToString(h[:])
}


