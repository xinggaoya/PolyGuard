package app

import (
	"PolyGuard/router"
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
	r.LoadHTMLGlob("templates/*")
	server := &http.Server{
		Addr:    ":6616",
		Handler: r,
	}

	router.InitRouter(r)
	//system.InitScanner()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
