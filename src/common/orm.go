package common

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var (
	db *gorm.DB
)

func InitOrm() {
	dsn := "root:root@tcp(localhost)/chun_search?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Failed connect to mysql")
	}
	db = gormDB
}

func DB() *gorm.DB {
	return db
}
