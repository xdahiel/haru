package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"haru/common"
	"net/http"
	"strings"
	"time"
)

type HaruClaim struct {
	Username string
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 60

func genToken(username string) (string, error) {
	if common.JwtSecret == "" {
		return "", fmt.Errorf("no jwt secret spetified")
	}

	access := HaruClaim{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "haru",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, access)
	return token.SignedString(common.JwtSecret)
}

func parseToken(tokenString string) (*HaruClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, new(HaruClaim),
		func(token *jwt.Token) (interface{}, error) {
			return common.JwtSecret, nil
		})

	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*HaruClaim); ok && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func jwtAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": "2003",
				"msg":  "请求头中的auth为空",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": "2004",
				"msg":  "请求头中的auth格式有误",
			})
			c.Abort()
			return
		}

		mc, err := parseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": "2005",
				"msg":  "无效的token",
			})
			c.Abort()
			return
		}

		c.Set("username", mc.Username)
		c.Next()
	}
}
