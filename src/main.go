package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"haru/common"
	"haru/middlewares"
)

func main() {
	common.Init()
	//user.Init()
	flag.Parse()

	r := gin.New()
	r.Use(middlewares.Middlewares()...)
}
