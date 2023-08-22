package system

import (
	"fmt"
	"net"
	"time"
)

/**
  @author: XingGao
  @date: 2023/8/22
**/

// InitScanner 定时扫描端口是否在监听
func InitScanner() {
	go func() {
		for {
			if isPortListening(6616) {
				fmt.Println("端口已经在监听")
			} else {
				fmt.Println("端口没有在监听")
			}
			time.Sleep(time.Second * 5)
		}
	}()
}

// IsPortListening 检查端口是否在监听
func isPortListening(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	err = listener.Close()
	if err != nil {
		return false
	}
	return false
}
