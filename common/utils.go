package common

import "golang.org/x/crypto/bcrypt"

func Init() {
	JwtSecret = generateSecreteRandomly()
	InitMysqlOrm()
	InitRDB()
}

func MD5(str string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hashedPassword)
}
