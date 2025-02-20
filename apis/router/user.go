package router

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/apis/controller"
)

func registerUserRoutes(rg *gin.RouterGroup) {
	rg.Group("/user/")
	rg.POST("register", controller.RegisterUser)
}
