package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"syncengin/pkg/core/memory"
	"syncengin/pkg/core/query"
	"syncengin/pkg/loader"
	"syncengin/pkg/log"
	"time"

	"github.com/gin-gonic/gin"
)

var config map[string]interface{}

func main() {
	start := time.Now()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	dir, err := os.Getwd()
	if err != nil {
		log.Error("Failed to get current working directory: " + err.Error())
		os.Exit(1)
	} else {
		memory.Append("dir", dir)
	}

	config = loader.LoadConfig(os.Args[1])
	router := gin.Default()
	routes(router)
	log.Success(fmt.Sprintf("Server started in %v", time.Since(start)))
	log.Info(fmt.Sprintf("Listening on port %v", config["port"]))
	log.Info(fmt.Sprintf("Serving files from %v", config["dir"]))
	router.Run(fmt.Sprintf(":%v", config["port"]))
}

func routes(router *gin.Engine) {
	router.GET("*path", handleRequest)
}

func handleRequest(c *gin.Context) {
	requestStart := time.Now()
	requestPath := c.Param("path")
	requestQuery := c.Request.URL.Query()

	log.Info("GET " + requestPath)

	if requestPath[len(requestPath)-1] == '/' {
		requestPath = requestPath[:len(requestPath)-1] + "/index.csv"
	} else {
		requestPath = requestPath + ".csv"
	}

	csvPath := path.Join(config["dir"].(string), requestPath)

	//check if path exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		log.Error("The requested endpoint " + c.Param("path") + " -> " + requestPath + " does not exist")
		c.JSON(http.StatusNotFound, gin.H{"error": "The requested endpoint does not exist"})
		return
	} else if err != nil {
		log.Error("An error occurred while checking the path: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	rows, err := loader.ReadCSV(csvPath)
	if err != nil {
		log.Error("An error occurred while reading the CSV file: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
	} else {
		response := query.Run(rows, requestQuery)
		if len(response) == 0 {
			response = []map[string]string{}
		}
		c.JSON(http.StatusOK, response)
	}

	log.Success(fmt.Sprintf("GET %v took %v", c.Param("path"), time.Since(requestStart)))
}
