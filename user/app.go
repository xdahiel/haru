package user

import (
	"github.com/gin-gonic/gin"
	"haru/user/controller"
	"haru/user/model"
)

func Init() {
	model.InitUser()
	model.InitAdvertise()
}

func InitRouter(v1 *gin.RouterGroup) {
	v1.Group("user")
	{
		v1.POST("/register", controller.Register)
		v1.POST("/login", controller.Login)
	}
}

func InitUserRouter(v1 *gin.RouterGroup) {
	v1.Group("user")
	{
		v1.POST("/upload", controller.Upload)
		v1.POST("/addAdv", controller.AddAdvertise)
	}
}
