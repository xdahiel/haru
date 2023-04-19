package common

import (
	"github.com/boltdb/bolt"
	"log"
)

var BoltDB *bolt.DB

func InitBoltDB() {
	gormDb, err := bolt.Open("haru.db", 0600, nil)
	if err != nil {
		log.Fatal("failed open/create local db:", err)
	}
	BoltDB = gormDb
}

func GetBoltDB() *bolt.DB {
	return BoltDB
}
