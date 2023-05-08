package common

import (
	"haru/logs"
	"math/rand"
)

var (
	JwtSecret string
)

func generateSecreteRandomly() string {
	dict := "0123456789qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	res := ""
	for i := 0; i < 25; i++ {
		p := rand.Int() % 62
		res += dict[p : p+1]
	}
	JwtSecret = res
	logs.Debug("random jwt secret: %v", JwtSecret)
	return res
}
