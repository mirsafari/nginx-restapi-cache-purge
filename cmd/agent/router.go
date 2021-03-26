package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	models "github.com/mirsafari/nginx-restapi-cache-purge/internal/models"
	validate "github.com/mirsafari/nginx-restapi-cache-purge/internal/validation"
	handlers "github.com/mirsafari/nginx/restapi-cache-purge/internal/handlers"
)

// Function that defines router groups. Currently we only have /v1 and are calling functions to add subroutes under v1
func getRoutes() {
	v1 := router.Group("/v1")

	addHealthceckRoute(v1)
	addCacheRoute(v1)
	handlers.AddSwaggerRoute(v1)
}

// Healthcheck route
func addHealthceckRoute(rg *gin.RouterGroup) {
	ping := rg.Group("/ping")
	a := models.PingEndpointResponse{}
	fmt.Println(a)
	ping.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
}

// Defines /v1/cache endpoint
func addCacheRoute(rg *gin.RouterGroup) {
	caches := rg.Group("/cache")
	// Handler for GET /v1/cache/:domain that check if cache for domain exists and returns cache size
	// User middleware checkContentType, checkAPIKey and checkDomainName on this call

	caches.GET("/:domain", validate.ContentType(), validate.APIKey(APIKey), validate.DomainName(), func(c *gin.Context) {
		a := models.CacheEndpointResponseGET{}
		fmt.Println(a)
		handlers.CacheGET(c, cachePath)
	})

	caches.DELETE("/:domain", validate.ContentType(), validate.APIKey(APIKey), validate.DomainName(), func(c *gin.Context) {
		a := models.CacheEndpointResponseDELETE{}
		fmt.Println(a)
		handlers.CacheDelete(c, cachePath)
	})
}
