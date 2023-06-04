package controller

import (
	"github.com/gin-gonic/gin"
	"haru/common"
	"haru/common/validator"
	"haru/logs"
	"haru/user/model"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	Role     string `json:"role" form:"role"`
}

func Register(c *gin.Context) {
	var rr RegisterRequest
	if err := c.ShouldBind(&rr); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "无法解析参数!",
		})
		logs.Info("failed resolve parameter: %v", err)
		return
	}
	logs.Info("username: %v, password: %v", rr.Username, rr.Password)

	if (!validator.ValidateUsername(rr.Username)) ||
		(!validator.ValidateEmail(rr.Email)) ||
		(!validator.ValidatePhone(rr.Phone)) ||
		(!validator.ValidatePassword(rr.Password)) {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数错误!",
		})
		logs.Error("Invalid parameter format.")
		return
	}

	u, err := model.FindUserByEmail(rr.Username)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2003",
			"msg":  "内部服务错误",
		})
		logs.Error("find user error: %v", err)
		return
	}

	if u != nil && len(u) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": "2004",
			"msg":  "用户已存在",
		})
		logs.Info("注册重名用户")
		return
	}

	err = model.AddUser(&model.User{
		ID:       0,
		Username: rr.Username,
		Email:    rr.Email,
		Phone:    rr.Phone,
		Password: common.MD5(rr.Password),
		Role:     rr.Role,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2003",
			"msg":  "内部服务错误",
		})
		logs.Error("add user error: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "2000",
		"msg":  "",
	})
}
