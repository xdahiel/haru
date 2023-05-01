package common

func Init() {
	JwtSecret = generateSecreteRandomly()
	InitMysqlOrm()
	InitRDB()
}
