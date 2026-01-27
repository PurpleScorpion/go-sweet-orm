package demo

import (
	"encoding/json"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
)

func registerDemo() {
	mapper.SetMySqlConf(mapper.MySQLConf{
		UserName: "root",
		Password: "root",
		DbName:   "demo",
		Port:     3308,
		Host:     "localhost",
	})
	mapper.RegisterMySql()
}

type user struct {
	Id       int    `json:"id" tableId:"id"`
	UserName string `json:"userName"`
	Age      int    `json:"age"`
	Word1    string `json:"word1"`
	Word2    string `json:"word2"`
	Word3    string `json:"word3"`
}

func (u user) TableName() string {
	return "user" // 返回正确的表名
}

func toString(data any) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}
