package main

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/common/middleware"
	"github.com/loveRyujin/go-mall/config"
	"net/http"
)

func main() {
	r := gin.New()
	r.Use(middleware.StartTrace())
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
	r.GET("/logger-test", func(ctx *gin.Context) {
		logger.New(ctx).Info("logger test", "key", "keyVal", "val", 2)
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	r.Run(":8080")
}
