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

	//2.获取旧链接
	oldConn, isExist := cm.connections[userId]

	//3.向map中添加新链接
	cm.connections[userId] = conn

	//4.释放锁
	cm.mu.Unlock()

	//5.如果链接已存在，关闭旧链接
	if isExist {

		//6.日志打印
		logrus.Warnf("[websocket-链接管理器] 存在旧链接: [%v]", userId)

		//7.关闭旧链接
		if err := oldConn.Close(); err != nil {
			logrus.Errorf("[websocket-链接管理器] 旧链接关闭异常: userId=[%s] err=[%v]", userId, err)
		} else {
			logrus.Debugf("[websocket-链接管理器] 旧链接关闭成功: [%v]", userId)
		}
	}

	//8.打印最终日志
	logrus.Debugf("[websocket-链接管理器] 添加链接成功: [%v]", userId)
}

// Remove 删除链接
func (cm *connManager) Remove(userId string) {

	//1.加锁
	cm.mu.Lock()
	defer cm.mu.Unlock()

	//2.从map中删除链接
	delete(cm.connections, userId)
	logrus.Debugf("[websocket-链接管理器] 删除链接成功: [%v]", userId)
}

// Close 关闭链接
func (cm *connManager) Close(userId string, exceptConn *websocket.Conn) {

	//1.加锁
	cm.mu.Lock()

	//2.获取链接
	conn, isExist := cm.connections[userId]

	//3.如果存在并且是同一个链接，移除链接
	if isExist && conn == exceptConn {
		delete(cm.connections, userId)
	}

	//4.释放锁
	cm.mu.Unlock()

	//5.如果存在并且是同一个链接
	if isExist && conn == exceptConn {

		//6.关闭链接
		if err := exceptConn.Close(); err != nil {
			logrus.Errorf("[websocket-链接管理器] 链接关闭异常: userId=[%s] err=[%v]", userId, err)
		}

		//7.打印日志
		logrus.Debugf("[websocket-链接管理器] 关闭链接成功: [%v]", userId)
	}
}

// CloseAll 关闭所有链接
func (cm *connManager) CloseAll() {

	//1.加锁
	cm.mu.Lock()

	//2.复制map
	copyConnections := cm.connections

	//3.重置map
	cm.connections = make(map[string]*websocket.Conn)

	//4.释放锁
	cm.mu.Unlock()

	//5.循环所有链接
	for userId, conn := range copyConnections {

		//6.关闭链接
		if err := conn.Close(); err != nil {
			logrus.Errorf("[websocket-链接管理器] 链接关闭异常: userId=[%s] err=[%v]", userId, err)
		}
	}

	//7.打印日志
	logrus.Infof("[websocket-链接管理器] 关闭所有链接成功")
}
