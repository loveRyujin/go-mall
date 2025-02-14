package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/common/app"
	"github.com/loveRyujin/go-mall/common/errcode"
	"github.com/loveRyujin/go-mall/logic/domainservice"
)

func AuthUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("my-go-mall")
		if len(token) != 64 { // 生成的token长度为64
			app.NewResponse(ctx).Error(errcode.ErrToken)
			return
		}
		tokenVerify, err := domainservice.NewUserDomainService(ctx).VerifyToken(token)
		if err != nil {
			app.NewResponse(ctx).Error(errcode.ErrServer)
			ctx.Abort()
			return
		}
		if tokenVerify != nil && !tokenVerify.Approved { // token验证未通过cl
			app.NewResponse(ctx).Error(errcode.ErrToken)
			ctx.Abort()
			return
		}
		ctx.Set("userId", tokenVerify.UserId)
		ctx.Set("sessionId", tokenVerify.SessionId)
		ctx.Next()
	}
}
