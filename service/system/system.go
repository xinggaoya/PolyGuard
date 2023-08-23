package system

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"syscall"
	"time"
)

/**
  @author: XingGao
  @date: 2023/8/22
**/

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Port    int    `json:"port"`
	IsAlive bool   `json:"is_alive"`
}

// InitScanner 定时扫描端口是否在监听
func InitScanner() {
	go func() {
		for {
			var serviceInfoVis []ServiceInfo
			serviceInfoVis = append(serviceInfoVis, ServiceInfo{
				Name:    "测试任务一",
				Path:    "C://Users/10322/Downloads/demo.java",
				Port:    8080,
				IsAlive: false,
			})
			for _, serviceInfo := range serviceInfoVis {
				if isPortInUse(serviceInfo.Port) {
					serviceInfo.IsAlive = true
				} else {
					// 重启服务
					serviceInfo.IsAlive = false
					RestartService(serviceInfo)
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

func isPortInUse(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return true // 端口已被占用
	}
	listener.Close()
	return false // 端口未被占用
}

// RestartService 重启服务
func RestartService(info ServiceInfo) {
	fmt.Println("重启服务")
	// 子进程执行命令
	cmd := exec.Command("java", "-java", info.Path)
	err := cmd.Start()
	if err != nil {
		fmt.Println("cmd.Start: ", err)
	}
	err = cmd.Wait()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// 子进程以非零状态退出
			status := exitErr.Sys().(syscall.WaitStatus)
			fmt.Printf("Child process exited with status %d\n", status.ExitStatus())
		}
	} else {
		fmt.Println("Child process exited successfully")
	}
}
