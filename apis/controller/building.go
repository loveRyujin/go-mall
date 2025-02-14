package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/apis/request"
	"github.com/loveRyujin/go-mall/common/app"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/config"
	"github.com/loveRyujin/go-mall/lib"
	"github.com/loveRyujin/go-mall/logic/appservice"
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
	service := appservice.NewDemoAppService(ctx)
	ids, err := service.GetDemoIds()
	if err != nil {
		app.NewResponse(ctx).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(ctx).Success(ids)
}

func TestCreateDemoOrder(ctx *gin.Context) {
	orderRequest := new(request.DemoOrderCreate)
	if err := ctx.ShouldBindJSON(orderRequest); err != nil {
		app.NewResponse(ctx).Error(errcode.ErrParams.WithCause(err))
		return
	}
	// 验证用户信息 Token 然后把UserID赋值上去 这里测试就直接赋值了
	orderRequest.UserId = 123453453
	service := appservice.NewDemoAppService(ctx)
	reply, err := service.CreateDemoOrder(orderRequest)
	if err != nil {
		app.NewResponse(ctx).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(ctx).Success(reply)
}

func TestHttpToolGet(ctx *gin.Context) {
	detail, err := lib.NewWhoisLib(ctx).GetHostIpDetail()
	if err != nil {
		app.NewResponse(ctx).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(ctx).Success(detail)
}

func TestHttpToolPost(ctx *gin.Context) {
	demoOrder, err := lib.NewDemoLib(ctx).TestPostCreateOrder()
	if err != nil {
		app.NewResponse(ctx).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(ctx).Success(demoOrder)
}

func TestMakeToken(ctx *gin.Context) {
	userService := appservice.NewUserAppService(ctx)
	token, err := userService.GenToken()
	if err != nil {
		if errors.Is(err, errcode.ErrUserInvalid) {
			logger.New(ctx).Error("invalid user is unable to generate token", "err", err)
			app.NewResponse(ctx).Error(errcode.ErrUserInvalid)
		} else {
			var appErr *errcode.AppError
			if ok := errors.As(err, appErr); ok {
				app.NewResponse(ctx).Error(appErr)
			}
		}
		return
	}
	app.NewResponse(ctx).Success(token)
}

func TestRefreshToken(ctx *gin.Context) {
	refreshToken := ctx.Query("refresh_token")
	if refreshToken == "" {
		app.NewResponse(ctx).Error(errcode.ErrParams)
		return
	}
	userService := appservice.NewUserAppService(ctx)
	token, err := userService.RefreshToken(refreshToken)
	if err != nil {
		if errors.Is(err, errcode.ErrTooManyRequests) {
			logger.New(ctx).Error("too many requests", "err", err)
			app.NewResponse(ctx).Error(errcode.ErrTooManyRequests)
		} else {
			var appErr *errcode.AppError
			if ok := errors.As(err, appErr); ok {
				app.NewResponse(ctx).Error(appErr)
			}
		}
		return
	}
	app.NewResponse(ctx).Success(token)
}

func TestAuthToken(ctx *gin.Context) {
	app.NewResponse(ctx).Success(gin.H{
		"user_id":    ctx.GetInt64("userId"),
		"session_id": ctx.GetString("sessionId"),
	})
	return
}
