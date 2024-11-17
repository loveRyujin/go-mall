package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/config"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("config-read", func(ctx *gin.Context) {
		database := config.Database
		ctx.JSON(http.StatusOK, gin.H{
			"type":     database.Type,
			"max_life": database.MaxLifeTime,
		})
	})
	r.Run(":10000")
}
