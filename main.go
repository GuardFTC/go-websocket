package main

import (
	"go-websocket/server"
	"go-websocket/ws"
)

func main() {

	//1.初始化ws链接管理器
	ws.InitConnManager()
	defer ws.ConnManager.CloseAll()

	//2.启动服务器
	server.StartServer()
}
