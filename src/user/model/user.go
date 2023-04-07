package model

import (
	"haru/common"
	"haru/logs"
	"log"
	"time"
)

type User struct {
	ID        int    `json:"id" gorm:"primaryKey;column:id;type:int;AUTO_INCREMENT"`
	Username  string `json:"username" gorm:"column:username;type:varchar(30)"`
	Password  string `json:"password" gorm:"column:password;type:varchar(255)"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64  `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func InitUser() {
	db := common.GetMysqlDB()
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalln("create user table failed")
	}
}

func AddUser(username, password string) error {
	db := common.GetMysqlDB()
	err := db.Model(&User{}).Create(User{
		Username:  username,
		Password:  password,
		CreatedAt: time.Now().Unix(),
	}).Error
	if err != nil {
		logs.Error("add user error: %v", err)
		return err
	}
	return nil
}

func FindUser(username string) ([]*User, error) {
	db := common.GetMysqlDB()
	var u []*User
	err := db.Model(new(User)).Where("username = ?", username).Find(&u).Error
	if err != nil {
		logs.Error("find user error: %v", err)
		return nil, err
	}
	return u, nil
}
