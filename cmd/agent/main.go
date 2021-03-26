package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

// Setup global variables
var router = gin.Default()
var cachePath string
var APIKey string

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
