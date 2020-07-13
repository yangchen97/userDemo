package dao

import (
	"context"
	"entryTask/model"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

const tokenExpire = 30 * time.Minute
const cacheExpire = 5 * time.Hour

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: 	  "127.0.0.1:6379",
		Password: "",
		DB: 	  0,
	})
}

func SetToken(token string) {
	err := rdb.Set(ctx, token, "a", tokenExpire).Err()
	if err != nil {
		log.Println(err)
	}
}

func Set(key, value string) {
	err := rdb.Set(ctx, key, value, cacheExpire).Err()
	if err != nil {
		log.Println(err)
	}
}

func Get(key string) string {
	val, _ := rdb.Get(ctx, key).Result()
	return val
}

func AddToken(token string, user *model.User) {
	err := rdb.Set(ctx, "user" + token, "1", tokenExpire).Err()
	if err != nil {
		log.Printf("set add token-user err: %v", err)
	}

	fmt.Printf(user.Username)
	err = rdb.HSet(ctx, token, user.ToMap()).Err()
	if err != nil {
		log.Printf("hset add token-user err: %v", err)
	}
}


func UpdateToken(token string) {
	err := rdb.Expire(ctx, "user" + token, tokenExpire).Err()
	if err != nil {
		log.Printf("update token expire time err: %v", err)
	}
}


func GetUserByToken(token string) model.User{
	m, err := rdb.HGetAll(ctx, token).Result()
	if err != nil {
		log.Printf("hgetall err: %v", err)
	}

	user := model.User{
		Username: m["username"],
		Nickname: m["nickname"],
		PicUrl:   m["pic_url"],
	}
	fmt.Println(user)
	return user

}


func UpdateUser(token string, user *model.User)  {
	err := rdb.Del(ctx, token).Err()
	if err != nil {
		log.Printf("del err: %v", err)
	}

	err = rdb.HSet(ctx, token, user.ToMap()).Err()
	if err != nil {
		log.Printf("del err: %v", err)
	}

}
