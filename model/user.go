package model

type User struct {
	Uid string `db:"uid"`
	Username string `db:"username"`
	Nickname string `db:"nickname"`
	Password string `db:"password"`
	PicUrl	string `db:"pic_url"`
}


func (user *User) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["username"] = user.Username
	m["nickname"] = user.Nickname
	m["pic_url"] = user.PicUrl
	return m
}