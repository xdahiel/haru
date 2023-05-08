package associate

import (
	"github.com/gin-gonic/gin"
	"haru/associate/controller"
)

func InitRouter(v1 *gin.RouterGroup) {
	v1.Group("asso")
	{
		//v1.POST("/add", controller.Add)
		v1.GET("/seek", controller.Seek)
	}
}
