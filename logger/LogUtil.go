package logger

import (
	"github.com/PurpleScorpion/go-sweet-orm/utils"
	"github.com/beego/beego/v2/core/logs"
	"os"
)

type LogUtil struct {
}

var log *logs.BeeLogger

func init() {

	runMode := os.Getenv("BEEGO_RUNMODE")
	filename := "go_ems" + ".log"
	var logFolderPath string
	if runMode == "prod" {
		// Read Configuration File
		// Define the folder path for storing logs
		logFolderPath = "/app/logs/"
	} else {
		// Define the folder path for storing logs
		logFolderPath = "logs/"
	}
	filename = logFolderPath + filename
	// Initialize log variable (10000 is the cache size)
	// logs.Async()
	// 创建一个日志器，可以给它指定一个名称，便于区分多个日志器
	log = logs.NewLogger()
	// 设置日志级别，例如：debug、info、warn、error、critical，默认为debug
	log.SetLevel(logs.LevelInfo)

	// 添加一个文件日志引擎，指定日志文件路径和模式（如按天分割、按大小分割等）

	js := utils.NewJSONObject()
	js.FluentPut("filename", filename)
	js.FluentPut("maxSize", 10*1024*1024)
	js.FluentPut("maxDays", 7)
	js.FluentPut("daily", true)
	js.FluentPut("maxBackups", 3)
	js.FluentPut("level", logs.LevelInfo)

	// 如果需要按小时分割文件，可以设置HourlyRolling为true
	// fileLogConfig.HourlyRolling = true
	log.Async()
	// 将文件日志引擎添加到日志器中
	if err := log.SetLogger(logs.AdapterFile, js.ToJsonString()); err != nil {
		panic("Failed to set file logger: " + err.Error())
	}
}

func Info(format string) {
	// 控制台打印
	logs.Info(format)
	// 文件记录
	log.Info(format)
}

func Warn(format string) {
	logs.Warn(format)
	log.Warn(format)
}

func Error(format string) {
	logs.Error(format)
	log.Error(format)
}
