package java

import (
	"PolyGuard/consts/taskConsts"
	"PolyGuard/service/db"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

// ServiceInfo 服务信息
type ServiceInfo struct {
	Id         string `json:"id" form:"id"`
	Name       string `json:"name" form:"name"`
	Path       string `json:"path" form:"path"`
	Port       int    `json:"port" form:"port"`
	IsAlive    bool   `json:"isAlive" form:"isAlive"`
	MaxRetries int    `json:"maxRetries" form:"maxRetries"`
	JvmOptions string `json:"jvmOptions" form:"jvmOptions"`
	Pid        int    `json:"pid" form:"pid"`
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
		pid := cmd.Process.Pid
		fmt.Printf("启动成功: %s, 进程ID: %d\n", task.Name, pid)
		task.IsAlive = true
		task.Pid = pid
		updateServiceStatus(task)

		err = cmd.Wait()
		if err != nil {
			retries++
			fmt.Printf("%s 进程异常关闭，尝试重启，重启次数: %d\n", task.Name, retries)
			time.Sleep(3 * time.Second) // 等待一段时间再重启

		} else {
			fmt.Println("进程退出: ", task.Name)
			task.IsAlive = false
			task.Pid = 0
			updateServiceStatus(task)
			break // 子进程成功退出，停止重启
		}
	}

	if retries == task.MaxRetries {
		fmt.Println(task.Name + " 已达到最大重启次数，不再重启")
		task.IsAlive = false
		task.Pid = 0
		updateServiceStatus(task)
	}
}

// StopService 停止jar包
func StopService(pid int) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	case "darwin", "linux":
		cmd = exec.Command("kill", "-9", strconv.Itoa(pid))
	default:
		return fmt.Errorf("不支持的操作系统")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	// 等待进程结束
	err = cmd.Wait()
	if err != nil {
		return err
	}
	fmt.Printf("进程已停止: %d\n", pid)
	return nil
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
