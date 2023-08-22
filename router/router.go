package router

import "github.com/gin-gonic/gin"

/**
  @author: XingGao
  @date: 2023/8/22
**/

// InitRouter 初始化路由
func InitRouter(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
}
