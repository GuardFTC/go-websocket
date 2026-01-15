// Package ws @Author:冯铁城 [17615007230@163.com] 2026-01-15 10:28:07
package ws

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// NormalHandler 创建 websocket 处理器
func NormalHandler(c *gin.Context) {

	//1.请求头获取用户ID
	userID := c.GetHeader("userID")
	if userID == "" {
		logrus.Warnf("[websocket-normal处理器] 用户ID为空")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//2.升级HTTP请求为WebSocket请求
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorf("[websocket-normal处理器] 升级HTTP请求为WebSocket请求失败: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	//3.将链接加入关闭器，并确保最终关闭链接
	ConnManager.Add(userID, conn)
	defer ConnManager.Close(userID)

	//4.定义消息处理逻辑，循环读取消息
	for {

		//5.读取消息
		messageType, p, err := conn.ReadMessage()
		if err != nil && !errors.Is(err, websocket.ErrCloseSent) {
			logrus.Errorf("[websocket-normal处理器] 读取消息失败: %v", err)
			break
		}

		//6.打印消息 模拟消息处理逻辑
		fmt.Printf("收到消息: %s\n", p)

		//7.推送消息给客户端，模拟服务端推送消息
		err = conn.WriteMessage(messageType, []byte("服务器已收到: "+string(p)))
		if err != nil {
			logrus.Errorf("[websocket-normal处理器] 推送消息给客户端失败: %v", err)
			break
		}
	}
}
