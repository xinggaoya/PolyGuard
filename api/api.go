package api

import (
	"PolyGuard/consts/taskConsts"
	"PolyGuard/service/db"
	"PolyGuard/service/java"
	task2 "PolyGuard/task"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
  @author: XingGao
  @date: 2023/8/22
**/

// Index 任务列表
func Index(c *gin.Context) {
	var list []java.ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	c.HTML(200, "index.html", list)
}

// AddPage 添加任务页面
func AddPage(c *gin.Context) {
	c.HTML(200, "add.html", nil)
}

// UpdatePage 修改任务页面
func UpdatePage(c *gin.Context) {
	var name = c.Query("name")
	var list []java.ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	var task java.ServiceInfo
	for _, v := range list {
		if v.Name == name {
			task = v
			break
		}
	}
	c.HTML(200, "update.html", task)
}

// Add 添加任务
func Add(c *gin.Context) {
	var task java.ServiceInfo
	err := c.ShouldBind(&task)
	if err != nil {
		fmt.Println(err)
	}
	var list []java.ServiceInfo
	err = db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	list = append(list, task)
	err = db.Set(taskConsts.TaskListKey, list)
	if err != nil {
		fmt.Println(err)
	}
	task2.InitTask()
	// 重定向
	c.Redirect(http.StatusFound, "/")
}

// Update 修改任务
func Update(c *gin.Context) {
	var task java.ServiceInfo
	// 接受表单数据
	err := c.ShouldBind(&task)
	if err != nil {
		fmt.Println(err)
	}
	var list []java.ServiceInfo
	err = db.Get(taskConsts.TaskListKey, &list)
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
	task2.InitTask()
	c.Redirect(http.StatusFound, "/")
}

// Delete 删除任务
func Delete(c *gin.Context) {
	name := c.Query("name")
	var list []java.ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	for i, v := range list {
		if v.Name == name {
			list = append(list[:i], list[i+1:]...)
			break
		}
	}
	err = db.Set(taskConsts.TaskListKey, list)
	if err != nil {
		fmt.Println(err)
	}
	c.Redirect(http.StatusFound, "/")
}

// Ping 测试接口
func Ping(c *gin.Context) {
	var taskList []java.ServiceInfo
	serviceInfos := append(taskList, java.ServiceInfo{Name: "测试", Path: "C://Users/10322/Downloads/demo.jar", Port: 8080, IsAlive: false})
	err := db.Set("task", serviceInfos)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
	})
}
