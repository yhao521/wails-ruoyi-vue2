package ylog

import (
	"fmt"
	"log"
	"mySparkler/pkg/constants"
	"os"
	"runtime"
	"time"

	"github.com/vrischmann/userdir"
)

// log 增加日志记录
func Log(content any) {
	fmt.Println("日志:", time.Now().Format(constants.TimeFormat), content)
	path := pathExist(fmt.Sprintf("%s/logs", getAppPath()))
	// 创建或打开日志文件
	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s.log", path, time.Now().Format(constants.DateFormat)), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	//defer logFile.Close()
	//记录文件路径和行号
	_, file, line, _ := runtime.Caller(1)
	// 初始化日志
	logger := log.New(logFile, "", log.LstdFlags)
	logger.Printf("\t文件路径：%s:%d\n日志内容：%s\n", file, line, content)
}

// pathExist 判断文件目录是否存在，不存在创建
func pathExist(path string) string {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	return path
}

// getAppPath 获取应用主目录
func getAppPath() string {
	dataDir := userdir.GetDataHome()
	return pathExist(fmt.Sprintf("%s/mySparklerFiles", dataDir))
}
