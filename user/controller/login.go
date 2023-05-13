package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"haru/logs"
	"haru/middlewares"
	"haru/user/model"
	"net/http"
)

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func Login(c *gin.Context) {
	var lr LoginRequest
	if err := c.ShouldBind(&lr); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "无法解析参数！",
		})
		logs.Info("failed resolve parameter: %v", err)
		return
	}

	u, err := model.FindUserByEmail(lr.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2002",
			"msg":  "内部服务错误",
		})
		logs.Info("Internal error: %v", err)
		return
	}

	if u == nil || len(u) != 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": "2003",
			"msg":  "用户不存在！",
		})
		logs.Info("Non-exist user.")
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u[0].Password), []byte(lr.Password)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2004",
			"msg":  "密码错误！",
		})
		logs.Info("Incorrect password: %v", err)
		return
	}

	token, _ := middlewares.GenToken(lr.Email)
	logs.Debug("get user: %v", fmt.Sprintf("%#v", u[0]))
	c.SetCookie("token", token, 3600*24, "/", "localhost", false, false)
	c.SetCookie("username", u[0].Username+u[0].Role+fmt.Sprintf("%02d", u[0].ID), 3600*24, "/", "localhost", false, false)
	c.SetCookie("email", u[0].Email, 3600*24, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{
		"code": "2000",
		"msg": gin.H{
			"token": token,
			"user":  u[0],
		},
	})
	return
}
