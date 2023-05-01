package common

import "math/rand"

var (
	JwtSecret string
)

func generateSecreteRandomly() string {
	dict := "0123456789qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	res := ""
	for i := 0; i < 256; i++ {
		p := rand.Int() % 62
		res += dict[p : p+1]
	}
	return res
}
