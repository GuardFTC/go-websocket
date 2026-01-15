// Package ws @Author:冯铁城 [17615007230@163.com] 2026-01-15 11:12:06
package ws

import (
	"sync"

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
	mu          sync.RWMutex
}

// newConnManager 创建链接管理器
func newConnManager() *connManager {
	return &connManager{
		connections: make(map[string]*websocket.Conn),
	}
}

// Get 获取链接
func (cm *connManager) Get(userId string) (*websocket.Conn, bool) {

	//1.加读锁
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	//2.获取链接，返回
	conn, isExist := cm.connections[userId]
	return conn, isExist
}

// Add 添加链接
func (cm *connManager) Add(userId string, conn *websocket.Conn) {

	//1.加锁
	cm.mu.Lock()
	defer cm.mu.Unlock()

	//2.如果链接已存在，关闭旧链接
	if oldConn, isExist := cm.connections[userId]; isExist {

		//3.日志打印
		logrus.Warnf("[websocket-链接管理器] 链接已存在: [%v]", userId)

		//4.关闭旧链接
		if err := oldConn.Close(); err != nil {
			logrus.Errorf("[websocket-链接管理器] 旧链接关闭异常: userId=[%s] err=[%v]", userId, err)
		}
	}

	//5.向map中添加新链接
	cm.connections[userId] = conn
	logrus.Infof("[websocket-链接管理器] 添加链接成功: [%v]", userId)
}

// Remove 删除链接
func (cm *connManager) Remove(userId string) {

	//1.加锁
	cm.mu.Lock()
	defer cm.mu.Unlock()

	//2.从map中删除链接
	delete(cm.connections, userId)
	logrus.Infof("[websocket-链接管理器] 删除链接成功: [%v]", userId)
}

// Close 关闭链接
func (cm *connManager) Close(userId string) {

	//1.加锁
	cm.mu.Lock()
	defer cm.mu.Unlock()

	//2.获取链接
	conn, isExist := cm.connections[userId]
	if !isExist {
		return
	}

	//3.关闭链接
	if err := conn.Close(); err != nil {
		logrus.Errorf("[websocket-链接管理器] 链接关闭异常: userId=[%s] err=[%v]", userId, err)
	}

	//4.移除链接
	delete(cm.connections, userId)

	//5.打印日志
	logrus.Infof("[websocket-链接管理器] 关闭链接成功: [%v]", userId)
}

// CloseAll 关闭所有链接
func (cm *connManager) CloseAll() {

	//1.加锁
	cm.mu.Lock()
	defer cm.mu.Unlock()

	//2.循环所有链接
	for userId, conn := range cm.connections {

		//3.关闭链接
		if err := conn.Close(); err != nil {
			logrus.Errorf("[websocket-链接管理器] 链接关闭异常: userId=[%s] err=[%v]", userId, err)
		}
	}

	//4.重置map
	cm.connections = make(map[string]*websocket.Conn)

	//5.打印日志
	logrus.Infof("[websocket-链接管理器] 关闭所有链接成功")
}
