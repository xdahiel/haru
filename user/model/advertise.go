package model

import (
	"fmt"
	"haru/common"
	"haru/logs"
	"log"
)

type Advertise struct {
	Id       int    `json:"id" gorm:"primaryKey;column:id;type:int;AUTO_INCREMENT"`
	Uid      int    `json:"uid" gorm:"colum:uid;type:int"`
	Username string `json:"username" gorm:"column:username;type:varchar(255)"`
	Keyword  string `json:"keyword" gorm:"colum:keyword;type:varchar(255)"`
	Handle   string `json:"handle" gorm:"column:handle;type:varchar(255)"`
	Link     string `json:"link" gorm:"column:link;type:varchar(255)"`
}

func InitAdvertise() {
	db := common.GetMysqlDB()
	err := db.AutoMigrate(&Advertise{})
	if err != nil {
		log.Fatalln("Failed create user table:", err)
	}
}

func AddAdvertise(adv *Advertise) error {
	logs.Debug("adv: %v", fmt.Sprintf("%#v", adv))
	db := common.GetMysqlDB()
	err := db.Model(&Advertise{}).Create(&adv).Error
	if err != nil {
		logs.Error("Failed add user: %v", err)
		return err
	}
	return nil
}

func FindAdvertiseByKeyword(keyword string) ([]*Advertise, error) {
	db := common.GetMysqlDB()
	var ads []*Advertise
	err := db.Model(new(Advertise)).Where("keyword = ?", keyword).Find(&ads).Error
	if err != nil {
		logs.Error("Fatal find advertise: %v", err)
		return nil, err
	}
	return ads, nil
}

func FindAdvertiseByUid(uid int) ([]*Advertise, error) {
	db := common.GetMysqlDB()
	var ads []*Advertise
	err := db.Model(new(Advertise)).Where("uid = ?", uid).Find(&ads).Error
	if err != nil {
		logs.Error("Fatal find advertise: %v", err)
		return nil, err
	}
	return ads, nil
}

func FindAdvertiseById(uid int) ([]*Advertise, error) {
	db := common.GetMysqlDB()
	var ads []*Advertise
	err := db.Model(new(Advertise)).Where("id = ?", uid).Find(&ads).Error
	if err != nil {
		logs.Error("Fatal find advertise: %v", err)
		return nil, err
	}
	return ads, nil
}

func DeleteAdvertiseById(uid int) error {
	db := common.GetMysqlDB()
	err := db.Model(new(Advertise)).Where("id = ?", uid).Delete(new(Advertise)).Error
	if err != nil {
		logs.Error("Fatal delete advertise: %v", err)
		return err
	}
	return nil
}
