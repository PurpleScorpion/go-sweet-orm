package logger

import (
	"github.com/beego/beego/v2/core/logs"
)

type LogUtil struct {
}

func Info(format string) {
	// 控制台打印
	logs.Info(format)
	// 文件记录
}

func Warn(format string) {
	logs.Warn(format)
}

func Error(format string) {
	logs.Error(format)
}
