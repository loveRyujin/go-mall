package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/common/app"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/config"
	"github.com/loveRyujin/go-mall/dal/dao"
)

func TestPing(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func TestConfigRead(ctx *gin.Context) {
	database := config.Database
	ctx.JSON(200, gin.H{
		"type":     database.Type,
		"max_life": database.Master.MaxLifeTime,
	})
}

func TestLogger(ctx *gin.Context) {
	logger.New(ctx).Info("logger test", "key", "keyVal", "val", 2)
	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}

func TestAccessLog(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status": "ok",
	})
}

func TestResponse(ctx *gin.Context) {
	data := map[string]int{
		"a": 1,
		"b": 2,
	}
	app.NewResponse(ctx).Success(data)
	return
}

func TestGormDbLogger(ctx *gin.Context) {
	demo := dao.NewDemoDao(ctx)
	records, err := demo.GetAllDemos()
	if err != nil {
		app.NewResponse(ctx).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(ctx).Success(records)
}
