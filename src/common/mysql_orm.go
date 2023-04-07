package common

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	mysqlDB *gorm.DB
)

func InitMysqlOrm() {
	dsn := "root:root@tcp(47.109.43.210:3306)/haru?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Failed connect to mysql")
	}
	mysqlDB = gormDB
}

func GetMysqlDB() *gorm.DB {
	return mysqlDB
}
