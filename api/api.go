package api

import (
	"PolyGuard/service/db"
	"github.com/gin-gonic/gin"
	"log"
)

/**
  @author: XingGao
  @date: 2023/8/22
**/

type Person struct {
	ID   int
	Name string
	Age  int
}

// Index 测试接口
func Index(c *gin.Context) {
	person := Person{ID: 123, Name: "Alice123", Age: 45}

	err := db.Set("key1", person)
	if err != nil {
		log.Fatal(err)
	}

	var retrievedPerson Person
	err = db.Get("key1", &retrievedPerson)
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(200, gin.H{
		"message": retrievedPerson,
	})
}
