package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"haru/common"
	"haru/logs"
	"haru/user"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags)

	common.Init()
	user.Init()
	flag.Parse()

	r := gin.New()
	v1 := r.Group("/api/v1")

	user.InitRouter(v1)

	//r.Use(middlewares.Middlewares()...)

	r.Static("/static", "web/public")
	r.LoadHTMLFiles("web/index.html", "web/login.html",
		"web/register.html", "web/result.html", "web/user.html")

	r.GET("/", index)
	r.GET("/register.html", register)

	err := r.Run(":8080")
	if err != nil {
		logs.Error("run gin error: %v", err)
		return
	}
}

func index(c *gin.Context) {
	logs.Info("visit index.")
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func register(c *gin.Context) {
	logs.Info("visit register.")
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
