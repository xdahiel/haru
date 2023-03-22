package model

import (
	"github.com/google/uuid"
	"haru/common"
	"log"
	"time"
)

type User struct {
	ID        string `json:"id" gorm:"primaryKey;column:id;type:varchar(255)"`
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

func Add(username, password string) error {
	db := common.GetMysqlDB()
	return db.Model(&User{}).Create(User{
		ID:        uuid.NewString(),
		Username:  username,
		Password:  password,
		CreatedAt: time.Now().Unix(),
	}).Error
}
