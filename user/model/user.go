package model

import (
	"haru/common"
	"haru/logs"
	"log"
)

type User struct {
	ID        int    `json:"id" gorm:"primaryKey;column:id;type:int;AUTO_INCREMENT"`
	Username  string `json:"username" gorm:"column:username;type:varchar(30)"`
	Email     string `json:"email" gorm:"column:email;type:varchar(255)"`
	Phone     string `json:"phone" gorm:"column:phone;type:varchar(25)"`
	Role      string `json:"role" gorm:"column:role;type:varchar(2)"`
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

func AddUser(usr *User) error {
	db := common.GetMysqlDB()
	err := db.Model(&User{}).Create(&usr).Error
	if err != nil {
		logs.Error("add user error: %v", err)
		return err
	}
	return nil
}

func FindUserByEmail(email string) ([]*User, error) {
	db := common.GetMysqlDB()
	var u []*User
	err := db.Model(new(User)).Where("email = ?", email).Find(&u).Error
	if err != nil {
		logs.Error("find user error: %v", err)
		return nil, err
	}
	return u, nil
}

func FindUserByID(id int) ([]*User, error) {
	db := common.GetMysqlDB()
	var u []*User
	err := db.Model(new(User)).Where("id = ?", id).Find(&u).Error
	if err != nil {
		logs.Error("find user error: %v", err)
		return nil, err
	}
	return u, nil
}
