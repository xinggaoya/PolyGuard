package app

import (
	"PolyGuard/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
  @author: XingGao
  @date: 2023/8/22
**/

func Run() {

	r := gin.New()
	r.LoadHTMLGlob("templates/*")
	server := &http.Server{
		Addr:    ":6616",
		Handler: r,
	}

	router.InitRouter(r)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
