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
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Index 测试接口
func Index(c *gin.Context) {
	var list []Person
	person := Person{ID: 123, Name: "Alice123", Age: 45}
	list = append(list, person)
	err := db.Set("key1", list)
	if err != nil {
		log.Fatal(err)
	}

	var retrievedPerson []Person
	err = db.Get("key1", &retrievedPerson)
	if err != nil {
		log.Fatal(err)
	}
	c.HTML(200, "index.html", retrievedPerson[0])
}
