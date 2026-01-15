package main

import (
	"go-websocket/server"
	"go-websocket/ws"
)

func main() {

	//1.初始化ws链接管理器
	ws.InitConnManager()

	//2.启动服务器，传入关闭回调
	server.StartServer(onShutdown)
}

// onShutdown 服务器关闭回调
func onShutdown() {

	//1.关闭所有webSocket链接
	ws.ConnManager.CloseAll()
}
