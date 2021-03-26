package main

import (
	"fmt"
	"net/http"
	handlers "nginx-restapi-cache-purge/pkg/handlers"
	models "nginx-restapi-cache-purge/pkg/models"
	validate "nginx-restapi-cache-purge/pkg/validation"

	"github.com/gin-gonic/gin"
	// handlers "github.com/mirsafari/nginx-restapi-cache-purge/pkg/handlers/"
	// models "github.com/mirsafari/nginx-restapi-cache-purge/pkg/models/"
	// validate "github.com/mirsafari/nginx-restapi-cache-purge/pkg/validation/"
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
	// User middleware checkContentType, checkAPIKey and checkDomainName on this call
	caches.GET("/:domain", validate.ContentType(), validate.APIKey(APIKey), validate.DomainName(), func(c *gin.Context) {
		a := models.CacheEndpointResponseGET{}
		fmt.Println(a)
		handlers.CacheGET(c, cachePath)
	})

	// User middleware checkContentType, checkAPIKey and checkDomainName on this call
	caches.DELETE("/:domain", validate.ContentType(), validate.APIKey(APIKey), validate.DomainName(), func(c *gin.Context) {
		b := models.CacheEndpointResponseDELETE{}
		fmt.Println(b)
		handlers.CacheDelete(c, cachePath)
	})
}
