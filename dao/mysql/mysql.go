package mysql

//var db *gorm.DB

import (
	"Team2048_Tiktok/conf"
	"Team2048_Tiktok/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

// Init 初始化MySQL连接
func Init(cfg *conf.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return
	}

	//自动建表，只需要启动系统就能在本地的MySQL中根据model中的结构体，按照tag的要求建立数据库表
	//在本地只能 启动一次 ，不需要重复建表
	DB.AutoMigrate(&model.User{}, &model.Video{}, &model.Comment{})

	return DB.DB().Ping()
}

// Close 关闭MySQL连接
func Close() {
	DB.Close()
}

/*


// Init 初始化MySQL连接
func Init(cfg *conf.MySQLConfig) (err error) {

}

// Close 关闭MySQL连接
func Close() {

}

*/
