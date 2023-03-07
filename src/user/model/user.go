package model

import (
	"haru/common"
	"github.com/google/uuid"
	"log"
	"time"
)

type User struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	CreatedTime time.Time `json:"created_time"`
}

func InitUser() {
	db := common.DB()
	err := db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalln("create user table failed")
	}
}

func Add(username, password string) error {
	db := common.DB()
	return db.Model(&User{}).Create(User{
		ID:          uuid.NewString(),
		Username:    username,
		Password:    password,
		CreatedTime: time.Now(),
	}).Error
}
