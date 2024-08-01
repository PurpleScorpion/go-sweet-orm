package logger

import (
	"fmt"
	"log"
)

type LogUtil struct {
}

func Info(format string, data ...interface{}) {
	// 控制台打印
	log.Println(fmt.Sprintf(format, data...))
}

func Warn(format string, data ...interface{}) {
	log.Println(fmt.Sprintf(format, data...))
}

func Error(format string, data ...interface{}) {
	log.Println(fmt.Sprintf(format, data...))
}
