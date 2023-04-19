package common

import "flag"

func Init() {
	flag.StringVar(&JwtSecret, "jwt", "haru", "random secret")
	InitMysqlOrm()
	InitRDB()
}
