package file

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/vrischmann/userdir"
)

func ExistDir(path string) {
	// 判断路径是否存在
	_, err := os.ReadDir(path)
	if err != nil {
		// 不存在就创建
		err = os.MkdirAll(path, fs.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// getAppPath 获取应用主目录
func GetAppPath() string {
	//获取系统我的文档目录
	dataDir := userdir.GetDataHome()
	//获取我的文档目录
	return PathExist(fmt.Sprintf("%s/mySparklerFiles", dataDir))
}

// pathExist 判断文件目录是否存在，不存在创建
func PathExist(path string) string {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path, os.ModePerm)
	}
	return path
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
