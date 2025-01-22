package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/apis/router"
	"github.com/loveRyujin/go-mall/common/enum"
	"github.com/loveRyujin/go-mall/config"
)

func main() {
	if config.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	router.RegisterRoutes(r)
	_ = r.Run(":10000")
}
