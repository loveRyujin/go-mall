package router

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/apis/controller"
)

func registerBuildingRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/building/")
	g.GET("ping", controller.TestPing)
	g.GET("config-read", controller.TestConfigRead)
	g.GET("logger-test", controller.TestLogger)
	g.POST("access-log-test", controller.TestAccessLog)
	g.GET("response-test", controller.TestResponse)
	g.GET("gorm-db-logger-test", controller.TestGormDbLogger)
}
