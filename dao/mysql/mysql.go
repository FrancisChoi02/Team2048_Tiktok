package mysql

//var db *gorm.DB

import (
	"Team2048_Tiktok/conf"
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
