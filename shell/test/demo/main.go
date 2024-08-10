package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// MyService 是我们自定义的类
type MyService struct {
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewMyService 创建并初始化 MyService
func NewMyService() *MyService {
	ctx, cancel := context.WithCancel(context.Background())
	return &MyService{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Run 开始执行 MyService 的主要工作
func (s *MyService) Run() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		// 模拟服务运行
		for {
			select {
			case <-s.ctx.Done():
				fmt.Println("Service is shutting down...")
				return
			default:
				fmt.Println("Service is running...")
				time.Sleep(1 * time.Second) // 模拟一些工作
			}
		}
	}()
}

// Stop 终止服务
func (s *MyService) Stop() {
	s.cancel()
	s.wg.Wait()
}

func main() {
	// 创建并运行服务
	service := NewMyService()
	service.Run()

	// 捕获系统中断信号（Ctrl+C）
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞，直到接收到中断信号
	select {
	case <-sigChan:
		fmt.Println("Received an interrupt, stopping service...")
		service.Stop()
	}

	fmt.Println("Service has stopped.")
}
