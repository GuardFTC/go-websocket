// Package server @Author:冯铁城 [17615007230@163.com] 2026-01-15 10:37:42
package server

import (
	"go-websocket/ws"

	"github.com/gin-gonic/gin"
)

// initRouter 初始化路由
func initRouter() *gin.Engine {

	//1.初始化路由
	router := gin.Default()

	//2.定义webSocket路由组
	wsGroup := router.Group("/ws")

	//3.定义上传日志路由
	wsGroup.GET("normal", ws.NormalHandler)

	//4.返回
	return router
}
