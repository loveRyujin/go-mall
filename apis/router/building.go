package router

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/apis/controller"
	"github.com/loveRyujin/go-mall/common/middleware"
)

func registerBuildingRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/building/")
	g.GET("ping", controller.TestPing)
	g.GET("config-read", controller.TestConfigRead)
	g.GET("logger-test", controller.TestLogger)
	g.POST("access-log-test", controller.TestAccessLog)
	g.GET("response-test", controller.TestResponse)
	g.GET("gorm-db-logger-test", controller.TestGormDbLogger)
	g.POST("create-demo-order", controller.TestCreateDemoOrder)
	g.GET("httptool-get-test", controller.TestHttpToolGet)
	g.GET("httptool-post-test", controller.TestHttpToolPost)
	// 测试token生成
	g.GET("token-make-test", controller.TestMakeToken)
	g.GET("token-refresh-test", controller.TestRefreshToken)
	g.GET("token-auth-test", middleware.AuthUser(), controller.TestAuthToken)
}
