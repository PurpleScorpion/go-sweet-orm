package demo

import (
	"encoding/json"
	"github.com/PurpleScorpion/go-sweet-orm/v3/mapper"
)

func registerDemo() {
	mapper.SetMySqlConf(mapper.MySQLConf{
		UserName: "root",
		Password: "MojiErEntIanXiAdi1gULiNgJinGguAi.520",
		DbName:   "demo",
		Port:     3308,
		Host:     "192.168.251.124",
	})
	mapper.RegisterMySql()
}

type user struct {
	Id       int    `json:"id" tableId:"id"`
	UserName string `json:"userName"`
	Age      int    `json:"age"`
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
