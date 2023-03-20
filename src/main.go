package main

import (
	"flag"
	"fmt"
	"haru/common"
)

func main() {
	common.Init()
	//user.Init()
	flag.Parse()
	fmt.Println(common.JwtSecret)
}
