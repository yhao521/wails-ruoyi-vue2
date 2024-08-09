package db

import (
	"mySparkler/config"
	"mySparkler/pkg/file"
	"mySparkler/pkg/mysql"
	"mySparkler/pkg/sqlite"

	"gorm.io/gorm"
)

// db 获取数据库操作对象和数据库初始化
func Db() *gorm.DB {
	return Dbp(file.GetAppPath())
}

// 初始化数据库
// func InitSchema(appPath string) {
// 	_ = Db(appPath).AutoMigrate(&model.Config{}, &system.SysConfig{})
// }

// db 获取数据库操作对象和数据库初始化
func Dbp(appPath string) *gorm.DB {
	driver := config.Database.Type
	switch driver {
	// case "redis":
	// 	return redisCache.NewRedisCache()
	case "mysql":
		return mysql.MysqlDb()
	case "sqlite":
		return sqlite.SqliteDb(appPath)
	}
	// return sqlite.SqliteDb("")
	//打开数据库，如果不存在，则创建
	db := sqlite.SqliteDb(appPath)
	return db
}
