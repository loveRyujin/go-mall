package router

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/common/middleware"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.Use(middleware.StartTrace(), middleware.LogAccess())
	routeGroup := engine.Group("")
	registerBuildingRoutes(routeGroup)
	registerOrderRoutes(routeGroup)
	registerUserRoutes(routeGroup)
}
