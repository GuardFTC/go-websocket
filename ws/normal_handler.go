// Package ws @Author:冯铁城 [17615007230@163.com] 2026-01-15 10:28:07
package ws

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// NormalHandler 创建 websocket 处理器
func NormalHandler(c *gin.Context) {

	//1.升级HTTP请求为WebSocket请求
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("升级失败:", err)
		return
	}
	defer conn.Close()

	//2.打印服务端已连接
	fmt.Println("服务端已连接")

	//3.定义消息处理逻辑，循环读取消息
	for {

		//4.读取消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取错误:", err)
			break
		}

		//5.打印消息 模拟消息处理逻辑
		fmt.Printf("收到消息: %s\n", p)

		//6.回显消息给客户端，模拟服务端推送消息
		err = conn.WriteMessage(messageType, []byte("服务器已收到: "+string(p)))
		if err != nil {
			log.Println("发送错误:", err)
			break
		}
	}
}
