package sqlite

import (
	"fmt"
	"log"
	"mySparkler/config"
	"os"
	"sync"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var once = sync.Once{}

type connect struct {
	client *gorm.DB
}

// 设置一个常量
var _connect *connect

func connectSqlite(appPath string) {

	//启用打印日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level: Silent、Error、Warn、Info
			Colorful:      false,       // 禁用彩色打印
		},
	)
	if config.Database.DbPath != "" {
		appPath = config.Database.DbPath
	}
	dbname := "mySparkler.db"
	if config.Database.DbFileName != "" {
		dbname = config.Database.DbFileName
	}
	dbFullPath := fmt.Sprintf("%s/"+dbname, appPath)
	// dsn := config.Database.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	dbs, errs := gorm.Open(sqlite.Open(dbFullPath), &gorm.Config{
		Logger: newLogger,
	})
	if errs != nil {
		panic("sqlite 数据库连接失败")
	}

	_connect = &connect{
		client: dbs,
	}
}

func SqliteDb(appPath string) *gorm.DB {
	if _connect == nil {
		once.Do(func() {
			connectSqlite(appPath)
		})
	}
	return _connect.client
}
