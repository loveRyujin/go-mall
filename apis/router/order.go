package router

import (
	"github.com/gin-gonic/gin"
	"github.com/loveRyujin/go-mall/apis/controller"
)

func registerOrderRoutes(rg *gin.RouterGroup) {
	g := rg.Group("/order/")
	// 创建订单
	g.POST("create", controller.OrderCreate)
	// 用户订单列表
	g.GET("user-order/", controller.UserOrders)
	// 订单详情
	g.GET(":order_id/info", controller.OrderInfo)
}
