package main

import (
	"fmt"
	"net/http"
	"statix/pkg/core/memory"
	"statix/pkg/loader"

	"github.com/gin-gonic/gin"
)

var config = loader.LoadConfig("config.json")

func main() {
	router := gin.Default()
	routes(router)
	router.Run(fmt.Sprintf(":%v", config["port"]))
}

func routes(router *gin.Engine) {
	fmt.Println(memory.Get("test")) //REMOVE
	router.GET("*path", handleRequest)
}

func handleRequest(c *gin.Context) {
	path := c.Param("path")
	c.JSON(http.StatusOK, gin.H{"path": path})
}
