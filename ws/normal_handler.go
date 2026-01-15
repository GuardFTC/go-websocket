// Package ws @Author:冯铁城 [17615007230@163.com] 2026-01-15 10:28:07
package ws

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	//3.将链接加管理器，并确保最终关闭链接
	ConnManager.Add(userID, conn)
	defer ConnManager.Close(userID, conn)

	//4.启动心跳协程
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go heartbeat(ctx, conn, userID)

	//5.定义消息处理逻辑，循环读取消息
	for {

		//6.设置读取超时（60秒内必须有网络活动，包括Pong）
		err = conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			logrus.Errorf("[websocket-normal处理器] 设置读取超时失败: userID=[%s] err=[%v]", userID, err)
			break
		}

		//7.读取消息
		messageType, p, err := conn.ReadMessage()

		//8.判断是否为正常关闭或服务端主动关闭
		if err != nil {
			if !(websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) ||
				websocket.IsUnexpectedCloseError(err) == false) {
				logrus.Errorf("[websocket-normal处理器] 读取消息失败: userID=[%s] err=[%v]", userID, err)
			}
			break
		}

		//9.打印消息 模拟消息处理逻辑
		fmt.Printf("收到消息: %s\n", p)

		//10.推送消息给客户端，模拟服务端推送消息
		err = conn.WriteMessage(messageType, []byte("服务器已收到: "+string(p)))
		if err != nil {
			logrus.Errorf("[websocket-normal处理器] 推送消息给客户端失败: %v", err)
			break
		}
	}
}

// heartbeat 心跳协程，定期发送Ping保持连接活跃
func heartbeat(ctx context.Context, conn *websocket.Conn, userID string) {

	//1.打印日志
	logrus.Debugf("[websocket-心跳] 开启心跳协程: userID=[%s]", userID)

	//2.创建定时器，每30秒触发一次
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		logrus.Debugf("[websocket-心跳] 关闭心跳协程: userID=[%s]", userID)
	}()

	//3.循环发送心跳
	for {
		select {

		//4.主协程退出，心跳协程也退出
		case <-ctx.Done():
			return

		//5.发送Ping帧（客户端会自动回复Pong）
		case <-ticker.C:
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
			if err != nil {
				logrus.Errorf("[websocket-心跳] 发送Ping失败: userID=[%s] err=[%v]", userID, err)
				return
			}
			logrus.Debugf("[websocket-心跳] 发送Ping成功: userID=[%s]", userID)
		}
	}
}
