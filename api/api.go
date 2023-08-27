package api

import (
	"PolyGuard/consts/taskConsts"
	"PolyGuard/service/db"
	"PolyGuard/service/java"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	var id = c.Query("id")
	var list []java.ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	var task java.ServiceInfo
	for _, v := range list {
		if v.Id == id {
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
	// UUID
	newUUID := uuid.New()
	task.Id = newUUID.String()
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
	go java.RunJar(task)
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
		if v.Id == task.Id {
			list[i] = task
			break
		}
	}
	err = db.Set(taskConsts.TaskListKey, list)
	if err != nil {
		fmt.Println(err)
	}
	java.StopService(task.Pid)
	go java.RunJar(task)
	c.Redirect(http.StatusFound, "/")
}

// Start 启动任务
func Start(c *gin.Context) {
	id := c.Query("id")
	var list []java.ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	for i, v := range list {
		if v.Id == id {
			go java.RunJar(list[i])
			break
		}
	}
	err = db.Set(taskConsts.TaskListKey, list)
	if err != nil {
		fmt.Println(err)
	}
	c.Redirect(http.StatusFound, "/")
}

// Stop 停止任务
func Stop(c *gin.Context) {
	id := c.Query("id")
	var list []java.ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	for i, v := range list {
		if v.Id == id {
			list[i].IsAlive = false
			list[i].Pid = 0
			java.StopService(v.Pid)
			break
		}
	}
	err = db.Set(taskConsts.TaskListKey, list)
	if err != nil {
		fmt.Println(err)
	}
	c.Redirect(http.StatusFound, "/")
}

// Delete 删除任务
func Delete(c *gin.Context) {
	id := c.Query("id")
	var list []java.ServiceInfo
	err := db.Get(taskConsts.TaskListKey, &list)
	if err != nil {
		fmt.Println(err)
	}
	for i, v := range list {
		if v.Id == id {
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
