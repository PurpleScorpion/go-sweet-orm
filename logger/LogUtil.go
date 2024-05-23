package logger

import (
	"log"
)

type LogUtil struct {
}

func Info(format string) {
	// 控制台打印
	log.Println(format)
	// 文件记录
}

func Warn(format string) {
	log.Println(format)
}

func Error(format string) {
	log.Println(format)
}
