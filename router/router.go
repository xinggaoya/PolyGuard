package router

import (
	"PolyGuard/api"
	"github.com/gin-gonic/gin"
)

/**
  @author: XingGao
  @date: 2023/8/22
**/

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) {
	r.GET("/", api.Index)
	// 添加任务页面
	r.GET("/add", api.AddPage)
	// 修改任务页面
	r.GET("/update", api.UpdatePage)

	// 添加任务
	r.POST("/add", api.Add)
	// 修改任务
	r.POST("/update", api.Update)
	// 删除任务
	r.GET("/delete", api.Delete)
	// ping
	r.GET("/ping", api.Ping)
}
