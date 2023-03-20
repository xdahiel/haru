package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"haru/common"
	"time"
)

type Access struct {
	Username string
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 60

func genToken(username string) (string, error) {
	access := Access{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "haru",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, access)
	return token.SignedString(common.JwtSecret)
}
