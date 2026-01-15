// Package server @Author:冯铁城 [17615007230@163.com] 2026-01-15 10:44:04
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// StartServer 启动服务
func StartServer(onShutdown func()) {

	//1.初始化路由
	router := initRouter()

	//2.创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: router,
	}

	//3.协程异步启动服务器
	go func() {
		logrus.Infof("[Server] 服务启动成功 监听端口: [%d]", 8080)
		if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			logrus.Fatalf("[Server] 启动服务器失败: [%v]", err)
		}
	}()

	//4.优雅关闭服务器
	waitForShutdown(server, onShutdown)
}

// waitForShutdown 优雅关闭服务器
func waitForShutdown(server *http.Server, onShutdown func()) {

	//1.创建信号通道
	quit := make(chan os.Signal, 1)

	//2.监听退出信号
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	//3.阻塞等待信号通道写入退出信号
	<-quit
	logrus.Infof("[Server] 接收到关闭信号，开始优雅关闭...")

	//4.执行关闭前的清理工作（如关闭所有 WebSocket 连接）
	if onShutdown != nil {
		onShutdown()
	}

	//5.设置关闭超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//6.关闭HTTP服务器，等待现有连接完成
	if err := server.Shutdown(ctx); err != nil {
		logrus.Errorf("[Server] 服务器关闭异常: [%v]", err)
	} else {
		logrus.Infof("[Server] 服务器关闭成功")
	}
}
