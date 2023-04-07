package user

import (
	"github.com/gin-gonic/gin"
	"haru/user/controller"
	"haru/user/model"
)

func Init() {
	model.InitUser()
}

func InitRouter(v1 *gin.RouterGroup) {
	v1.Group("user")
	{
		v1.POST("/register", controller.Register)
	}
}
