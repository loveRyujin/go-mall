package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/apis/request"
	"github.com/loveRyujin/go-mall/common/app"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/common/logger"
	"github.com/loveRyujin/go-mall/common/utils"
	"github.com/loveRyujin/go-mall/logic/appservice"
)

func RegisterUser(ctx *gin.Context) {
	registerReq := new(request.UserRegister)
	if err := ctx.ShouldBind(registerReq); err != nil {
		app.NewResponse(ctx).Error(errcode.ErrParams.WithCause(err))
		return
	}
	if !utils.VerifyPasswordComplexity(registerReq.Password) {
		logger.New(ctx).Warn("register user password complexity check failed")
		app.NewResponse(ctx).Error(errcode.ErrParams)
		return
	}
	// 注册用户
	userAppSvr := appservice.NewUserAppService(ctx)
	if err := userAppSvr.RegisterUser(registerReq); err != nil {
		if errors.Is(err, errcode.ErrUserNameOccupied) {
			app.NewResponse(ctx).Error(errcode.ErrUserNameOccupied)
		} else {
			app.NewResponse(ctx).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(ctx).SuccessOk()
	return
}
