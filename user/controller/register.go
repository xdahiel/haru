package controller

import (
	"github.com/gin-gonic/gin"
	"haru/logs"
	"haru/user/model"
	"log"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func Register(c *gin.Context) {
	var rr RegisterRequest
	if err := c.ShouldBind(&rr); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "无法解析参数!",
		})
		log.Printf("register error: %v", err)
		return
	}
	logs.Info("username: %v, password: %v", rr.Username, rr.Password)

	if len(rr.Username) < 6 || len(rr.Username) > 20 {
		c.JSON(http.StatusOK, gin.H{
			"code": "2002",
			"msg":  "用户名应在6~20位之间",
		})
		logs.Info("用户输入不规范")
		return
	}

	if len(rr.Password) < 6 || len(rr.Password) > 20 {
		c.JSON(http.StatusOK, gin.H{
			"code": "2002",
			"msg":  "密码应在6~20位之间",
		})
		logs.Info("密码输入不规范")
		return
	}

	u, err := model.FindUser(rr.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2003",
			"msg":  "内部服务错误",
		})
		logs.Error("find user error: %v", err)
		return
	}

	if u != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2004",
			"msg":  "用户已存在",
		})
		logs.Info("注册重名用户")
		return
	}

	err = model.AddUser(rr.Username, rr.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2003",
			"msg":  "内部服务错误",
		})
		logs.Error("add user error: %v", err)
		return
	}
}
