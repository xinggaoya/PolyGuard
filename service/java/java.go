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
	IsAutoRun  bool   `json:"isAutoRun" form:"isAutoRun"` // 是否自动启动
}

// RunJar 执行jar包
func RunJar(task ServiceInfo) {
	task.IsAutoRun = true
	retries := 0
	var list []ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	taskNew := ServiceInfo{}
	for _, v := range list {
		if v.Id == task.Id {
			taskNew = v
			break
		}
	}

	for retries < taskNew.MaxRetries && taskNew.IsAutoRun {
		if taskNew.JvmOptions != "" {
			taskNew.Path = taskNew.Path + " " + taskNew.JvmOptions
		}
		cmd := exec.Command("java", "-jar", taskNew.Path)
		err = cmd.Start()
		if err != nil {
			fmt.Println("启动失败:" + taskNew.Name)
		}
		pid := cmd.Process.Pid
		fmt.Printf("启动成功: %s, 进程ID: %d\n", taskNew.Name, pid)
		taskNew.IsAlive = true
		taskNew.Pid = pid
		updateServiceStatus(taskNew)

		err = cmd.Wait()
		if err != nil {
			retries++
			fmt.Printf("%s 进程异常关闭，尝试重启，重启次数: %d\n", task.Name, retries)
			time.Sleep(3 * time.Second) // 等待一段时间再重启
		} else {
			fmt.Println("进程退出: ", task.Name)
			task.IsAlive = false
			task.Pid = 0
			task.IsAutoRun = true
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
func StopService(task ServiceInfo) error {
	var cmd *exec.Cmd

	// 排除0
	if task.Pid == 0 {
		return fmt.Errorf("进程不存在")
	}

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(task.Pid))
	case "darwin", "linux":
		cmd = exec.Command("kill", "-9", strconv.Itoa(task.Pid))
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
	task.IsAutoRun = false
	updateServiceStatus(task)
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
