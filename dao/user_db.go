package dao


import (
	sql "database/sql"
	"entryTask/model"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)


var db *sql.DB


func init() {
	var err error
	db, err = sql.Open("mysql", "root:yangchen@tcp(localhost:3306)/yangchen?charset=utf8")
	if err != nil {
		log.Printf("database connection err:%v", err)
	}

	db.SetMaxIdleConns(2000)
	db.SetMaxOpenConns(1000)
	db.Ping()
}


func GetUserByUsername(username string) (*model.User, error) {
	stmt, err := db.Prepare("SELECT * FROM table_user WHERE username=? LIMIT 1")
	if err != nil {
		log.Printf("prepare statement getuser err:%v", err)
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)
	if row == nil {
		return nil, errors.New("username not found")
	}
	var user model.User
	row.Scan(&user.Uid, &user.Username, &user.Nickname, &user.Password, &user.PicUrl)

	return &user, nil
}


func UpdateUserNicknameByUsername(username string, newNickname string) bool {
	stmt, err := db.Prepare("UPDATE table_user SET nickname = ? WHERE username = ?")
	if err != nil {
		log.Printf("nickname prepare statement updateuser err:%v", err)
		return false
	}
	defer stmt.Close()
	res, err := stmt.Exec(newNickname, username)
	if err != nil {
		log.Printf("nickname exec statement updateuser err:%v", err)
		return false
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("nickname updateuser err:%v", err)
		return false
	}

	if cnt == 0 {
		log.Printf("nickname 0 row affected err:%v", err)
		return false
	}


	return true


}



func UpdateUserPicUrlByUsername(username string, newPicUrl string) bool {
	stmt, err := db.Prepare("UPDATE table_user SET pic_url = ? WHERE username = ?")
	if err != nil {
		log.Printf("pic_url prepare statement updateuser err:%v", err)
		return false
	}
	defer stmt.Close()
	fmt.Println(newPicUrl)
	res, err := stmt.Exec(newPicUrl, username)
	if err != nil {
		log.Printf("pic_url exec statement updateuser err:%v", err)
		return false
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		log.Printf("pic_url updateuser err:%v", err)
		return false
	}

	if cnt == 0 {
		log.Printf("pic_url 0 row affected err:%v", err)
		return false
	}


	return true


}