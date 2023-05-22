package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"haru/common"
	"haru/common/validator"
	"haru/logs"
	"haru/user/model"
	"net/http"
	"strconv"
)

func User(c *gin.Context) {
	sid := c.Query("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "用户不存在！",
		})
		logs.Error("Invalid input: %v", err)
		return
	}

	user, err := model.FindUserByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "内部错误",
		})
		logs.Error("Failed found user: %v", err)
		return
	}

	if user == nil || len(user) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "用户不存在",
		})
		logs.Error("Get: %v", user)
		return
	}

	c.HTML(http.StatusOK, "user.html", gin.H{
		"code": "2000",
		"user": user[0],
	})

}

type UpdateUserRequest struct {
	ID       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
	Phone    string `json:"phone" form:"phone"`
}

func UpdateUserInfo(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数错误!",
		})
		logs.Error("Invalid parameter: %v", err)
		return
	}

	if !validator.ValidateUsername(req.Username) ||
		!validator.ValidatePhone(req.Phone) {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数错误!",
		})
		logs.Error("Invalid parameter format.")
		return
	}

	user := &model.User{
		ID:       req.ID,
		Username: req.Username,
		Phone:    req.Phone,
	}

	err := model.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "更新信息失败！",
		})
		logs.Error("Failed update user info: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "2000",
		"msg":  "",
	})
}

type UpdatePasswordRequest struct {
	ID          int    `json:"id" form:"id"`
	OldPassword string `json:"oldPassword" form:"oldPassword"`
	NewPassword string `json:"newPassword" form:"newPassword"`
}

func UpdateUserPassword(c *gin.Context) {
	var req UpdatePasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数错误!",
		})
		logs.Error("Invalid parameter: %v", err)
		return
	}

	if !validator.ValidatePassword(req.OldPassword) ||
		!validator.ValidatePassword(req.NewPassword) {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数错误!",
		})
		logs.Error("Invalid parameter format.")
		return
	}

	u, err := model.FindUserByID(req.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "内部错误!",
		})
		logs.Error("Failed found user: %v", err)
		return
	}

	if u == nil || len(u) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "内部错误!",
		})
		logs.Error("Failed found user.")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u[0].Password), []byte(req.OldPassword)) != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "原密码不正确!",
		})
		logs.Error("Invalid Parameter: incorrect old password")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u[0].Password), []byte(req.NewPassword)) == nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "旧密码和新密码一致!",
		})
		logs.Error("old password is equal to new password")
		return
	}

	hashedNewPassword := common.MD5(req.NewPassword)
	user := &model.User{
		ID:       req.ID,
		Password: hashedNewPassword,
	}

	err = model.UpdateUserPassword(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "更新信息失败！",
		})
		logs.Error("Failed update user info: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "2000",
		"msg":  "",
	})
}
