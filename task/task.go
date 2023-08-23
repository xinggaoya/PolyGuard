package task

import (
	"PolyGuard/service/db"
	"PolyGuard/service/java"
	"fmt"
	"net"
)

// InitTask 初始化任务
func InitTask() {
	var taskList []java.ServiceInfo
	db.Get("task", &taskList)
	for _, task := range taskList {
		if !isPortInUse(task.Port) {
			// 重启服务
			go java.RunJar(task)
		}
		task.IsAlive = true
	}
}

// 检查是否使用
func isPortInUse(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return true // 端口已被占用
	}
	listener.Close()
	return false // 端口未被占用
}
