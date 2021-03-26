package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()
var cachePath string
var APIKey string

// @title nginx-restapi-cache-purge
// @version 0.0.1
// @description A simple API for deleting Nginx cache via HTTP calls. It enables you to get status and size of cache for certian domain and also delete cache for that domain

// @contact.email miroslav.safaric@gmail.com

// @host localhost:5000
// @BasePath /

// Run will start the server
func main() {

	flag.StringVar(&cachePath, "cache-path", os.Getenv("NGINX_CACHE_DIR"), "Path to folder containing Nginx cahce")
	flag.StringVar(&APIKey, "api-key", os.Getenv("APIKEY"), "API Key used to access this service (min. 16 chars)")
	flag.Parse()

	// Simple validation for input flags
	if len(cachePath) == 0 && len(APIKey) < 16 {
		fmt.Println("Please provide valid API key and Nginx cache dir")
		os.Exit(1)
	}

	// getRoutes will create our routes of our entire application
	// this way every group of routes can be defined in their own file
	// so this one won't be so messy
	getRoutes()

	// Allow gin to return 405 if method is not defined in router
	router.HandleMethodNotAllowed = true

	// Start server
	router.Run(":5000")
}
