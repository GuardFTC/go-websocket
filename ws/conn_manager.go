// Package ws @Author:冯铁城 [17615007230@163.com] 2026-01-15 11:12:06
package ws

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// ConnManager 链接管理器
var ConnManager *connManager

// InitConnManager 初始化链接管理器
func InitConnManager() {
	ConnManager = newConnManager()
	logrus.Infof("[websocket-链接管理器] 链接管理器初始化成功")
}

// connManager 链接管理器
type connManager struct {
	connections map[string]*websocket.Conn
}

// newConnManager 创建链接管理器
func newConnManager() *connManager {
	return &connManager{
		connections: make(map[string]*websocket.Conn),
	}
}

// Add 添加链接
func (cm *connManager) Add(userId string, conn *websocket.Conn) {
	cm.connections[userId] = conn
	logrus.Infof("[websocket-链接管理器] 添加链接成功: [%v]", userId)
}

// Get 获取链接
func (cm *connManager) Get(userId string) (*websocket.Conn, bool) {
	conn, isExist := cm.connections[userId]
	return conn, isExist
}

// Remove 删除链接
func (cm *connManager) Remove(userId string) {
	delete(cm.connections, userId)
	logrus.Infof("[websocket-链接管理器] 删除链接成功: [%v]", userId)
}

// CloseAll 关闭所有链接
func (cm *connManager) CloseAll() {

	//1.循环所有链接
	for userId, conn := range cm.connections {

		//2.关闭链接
		if err := conn.Close(); err != nil {
			logrus.Errorf("[websocket-链接管理器] 链接关闭异常: [%v]", err)
		}

		//3.移除链接
		cm.Remove(userId)
	}

	//4.打印日志
	logrus.Infof("[websocket-链接管理器] 关闭所有链接成功")
}

// Close 关闭链接
func (cm *connManager) Close(userId string) {

	//1.获取链接
	conn, isExist := cm.Get(userId)
	if !isExist {
		return
	}

	//2.关闭链接
	if err := conn.Close(); err != nil {
		logrus.Errorf("[websocket-链接管理器] 关闭链接异常: [%v]", err)
	}

	//3.移除链接
	cm.Remove(userId)

	//4.打印日志
	logrus.Infof("[websocket-链接管理器] 关闭链接成功: [%v]", userId)
}
