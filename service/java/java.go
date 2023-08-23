package java

import (
	"PolyGuard/consts/taskConsts"
	"PolyGuard/service/db"
	"fmt"
	"os/exec"
	"time"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	Name       string `json:"name" form:"name"`
	Path       string `json:"path" form:"path"`
	Port       int    `json:"port" form:"port"`
	IsAlive    bool   `json:"isAlive" form:"isAlive"`
	MaxRetries int    `json:"maxRetries" form:"maxRetries"`
	JvmOptions string `json:"jvmOptions" form:"jvmOptions"`
}

// RunJar 执行jar包
func RunJar(task ServiceInfo) {
	retries := 0
	for retries < task.MaxRetries {
		if task.JvmOptions != "" {
			task.Path = task.Path + " " + task.JvmOptions
		}
		cmd := exec.Command("java", "-jar", task.Path)
		err := cmd.Start()
		if err != nil {
			fmt.Println("启动失败:" + task.Name)
		}
		fmt.Println("启动成功:" + task.Name)
		task.IsAlive = true
		updateServiceStatus(task)
		err = cmd.Wait()
		if err != nil {
			fmt.Println(task.Name + " 进程异常关闭，尝试重启")
			retries++
			fmt.Printf("重启次数：%d\n", retries)
			time.Sleep(time.Second) // 等待一段时间再重启
		} else {
			fmt.Println("进程退出: ", task.Name)
			break // 子进程成功退出，停止重启
		}
		task.IsAlive = false
		updateServiceStatus(task)
	}

	if retries == task.MaxRetries {
		fmt.Println(task.Name + " 已达到最大重启次数，不再重启")
	}
}

// 修改服务状态
func updateServiceStatus(task ServiceInfo) {
	var list []ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	for i, v := range list {
		if v.Name == task.Name {
			list[i] = task
			break
		}
	}
	err = db.Set(taskConsts.TaskListKey, list)
	if err != nil {
		fmt.Println(err)
	}
}
