package main

import (
	"autobase/pkg/core/memory"
	"autobase/pkg/core/query"
	"autobase/pkg/loader"
	"autobase/pkg/log"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var config map[string]interface{}
var keys map[string]string

type Response struct {
	Status string              `json:"status"`
	Count  int                 `json:"count,omitempty"`
	Data   []map[string]string `json:"data,omitempty"`
	Error  string              `json:"error,omitempty"`
}

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

	if len(os.Args) < 2 {
		log.Error("No config file provided")
		os.Exit(1)
	}

	config = loader.LoadConfig(os.Args[1])
	loadAPIKeys()
	router := gin.Default()
	router.Use(keyAuthMiddleware())
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

	response := query.Run(csvPath, requestQuery, c.GetString("access_level"))

	if response.Status == "error" {
		c.JSON(http.StatusInternalServerError, response)
	} else if response.Status == "ClientError" {
		c.JSON(http.StatusBadRequest, response)
	} else {
		c.JSON(http.StatusOK, response)
		log.Success(fmt.Sprintf("GET %v took %v", c.Param("path"), time.Since(requestStart)))
	}
}

func loadAPIKeys() {
	apiKeysJson := os.Getenv("AB_AUTH")
	if apiKeysJson != "" {
		err := json.Unmarshal([]byte(apiKeysJson), &keys)
		if err != nil {
			log.Error("Failed to parse AB_AUTH: " + err.Error())
		}

		for key, value := range keys {
			if value != "read" && value != "write" {
				log.Warning("Invalid access level for key '" + key + "'. Deleting key")
				delete(keys, key)
			}
		}
	}
}

func keyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if len(keys) == 0 {
			c.Next()
			return
		}

		//find length of read and write keys
		readKeys := 0
		writeKeys := 0
		for _, value := range keys {
			if value == "read" {
				readKeys++
			} else if value == "write" {
				writeKeys++
			}
		}

		authHeader := c.GetHeader("Authorization")
		if readKeys > 0 && writeKeys > 0 && authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// if no read keys specified, it means read doesn't require auth
		if readKeys == 0 && writeKeys > 0 && authHeader == "" {
			c.Set("access_level", "read")
			c.Next()
			return
		}

		// if no write keys specified, it means write doesn't require auth
		if writeKeys == 0 && readKeys > 0 && authHeader == "" {
			c.Set("access_level", "write")
			c.Next()
			return
		}

		// Split the auth header to get the token part
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		token := parts[1]
		accessLevel, exists := keys[token]
		if !exists {
			c.JSON(403, gin.H{"error": "Invalid auth key"})
			c.Abort()
			return
		}

		c.Set("access_level", accessLevel)
		c.Next()
	}
}
