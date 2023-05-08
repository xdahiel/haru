package controller

import (
	"github.com/gin-gonic/gin"
	"haru/associate/types"
	"net/http"
)

//type AddRequest struct {
//	Sentence string `json:"sentence" form:"sentence"`
//}

//func Add(c *gin.Context) {
//	var ar AddRequest
//	if err := c.ShouldBind(&ar); err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": "2001",
//			"msg":  "参数错误！",
//		})
//		logs.Error("Failed binding parameter: %v", err)
//		return
//	}
//
//	types.GetTrie().Insert(ar.Sentence)
//	c.JSON(http.StatusOK, gin.H{
//		"code": "2000",
//		"msg":  "success",
//	})
//}

func Seek(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": "2002",
			"msg":  []string{},
		})
		return
	}
	res := types.GetTrie().Seek(query)
	if len(res) > 6 {
		res = res[:6]
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "2002",
		"msg":  res,
	})
}
