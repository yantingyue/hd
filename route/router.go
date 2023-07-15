package route

import (
	"demo/hd/handlers"
	"github.com/gin-gonic/gin"
)

func Router(g *gin.Engine) {
	// 抢单
	g.POST("/api/grab/order", handlers.HdHandler.GrabOrder)
	// 停止抢单
	g.POST("/api/grab/order/stop", handlers.HdHandler.StopGrabOrder)
	// 抢单任务列表
	g.POST("/api/grab/order/task/list", handlers.HdHandler.GrabOrderTaskList)
	// 合成抽签
	g.POST("/api/combine/draw")
}
