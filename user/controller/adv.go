package controller

import (
	"github.com/gin-gonic/gin"
	"haru/logs"
	"haru/user/model"
	"net/http"
	"strconv"
)

func Company(c *gin.Context) {
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

	advs, err := model.FindAdvertiseByUid(user[0].ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "内部错误！",
		})
		logs.Error("Failed found advs by uid: %v", err)
		return
	}

	c.HTML(http.StatusOK, "company.html", gin.H{
		"code": "2000",
		"user": user[0],
		"advs": advs,
	})

}

type AddAdvertiseRequest struct {
	Id       string `json:"id" form:"id"`
	Keyword  string `json:"keyword" form:"keyword"`
	Handle   string `json:"handle" form:"handle"`
	Link     string `json:"link" form:"link"`
	Username string `json:"username" form:"username"`
}

func AddAdvertise(c *gin.Context) {
	var req AddAdvertiseRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数解析错误！",
		})
		logs.Error("Invalid input: %v", err)
		return
	}

	id, err := strconv.Atoi(req.Id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数解析错误！",
		})
		logs.Error("Invalid input: %v", err)
		return
	}

	adv := model.Advertise{
		Id:       0,
		Uid:      id,
		Keyword:  req.Keyword,
		Handle:   req.Handle,
		Link:     req.Link,
		Username: req.Username,
	}
	err = model.AddAdvertise(&adv)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "内部错误！",
		})
		logs.Error("Failed add advertise", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "2000",
		"msg":  "",
	})
}

func DeleteAdv(c *gin.Context) {
	sid := c.Query("id")
	id, err := strconv.Atoi(sid)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "参数不合法!",
		})
		logs.Error("Invalid parameter: %v", err)
		return
	}

	cur, err := model.FindAdvertiseById(id)
	if err != nil || cur == nil || len(cur) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "该条信息不存在！",
		})
		logs.Error("Failed found advertise: %v", err)
		return
	}

	err = model.DeleteAdvertiseById(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": "2001",
			"msg":  "删除失败!",
		})
		logs.Error("Failed delete advertise: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "2000",
		"msg":  "",
	})
}
