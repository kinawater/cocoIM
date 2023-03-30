package logger

import (
	"cocoIM/config"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

var LogFormatTime = "20060102"

// openLogFile 打开日志文件
func openLogFile() *os.File {
	// 文件名 +  时间 + 扩展名
	fileName := fmt.Sprintf("%s%s.%s", config.LoggerConf.LogSaveName, time.Now().Format(LogFormatTime), config.LoggerConf.LogFileExt)
	filePath := config.LoggerConf.LogPath + fileName
	_, err := os.Stat(filePath)
	switch {
	case errors.Is(err, fs.ErrExist):
		// 创建文件
		mkLogDir()
	case errors.Is(err, fs.ErrPermission):
		// 没权限
		log.Fatalf("日志文件写入无权限:%v", err)
	}
	handleFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("日志文件打开失败:%v", err)
	}
	return handleFile
}
func mkLogDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+config.LoggerConf.LogPath, os.ModePerm)
	if err != nil {
		log.Fatalf("日志文件创建失败:%v", err)
	}
}
