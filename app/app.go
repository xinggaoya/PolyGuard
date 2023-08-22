package app

import (
	"PolyGuard/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

/**
  @author: XingGao
  @date: 2023/8/22
**/

func Run() {

	r := gin.New()
	fmt.Println(gin.Version)

	server := &http.Server{
		Addr:    ":6616",
		Handler: r,
	}

	router.InitRouter(r)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
