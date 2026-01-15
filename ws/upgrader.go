// Package ws @Author:冯铁城 [17615007230@163.com] 2026-01-15 10:28:07
package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// upgrader 创建 websocket 升级器对象，供包内所有 handler 使用
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
