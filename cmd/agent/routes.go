package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	validate "github.com/mirsafari/nginx-restapi-cache-purge/pkg/validation"
)

// Function that defines router groups. Currently we only have /v1 and are calling functions to add subroutes under v1
func getRoutes() {
	v1 := router.Group("/v1")
	addHealthceckRoute(v1)
	addCacheRoute(v1)
}

// Simple handler for /v1/ping route that returns success. Since this is only a healthcheck endpoint, we dont put it in another file
// HealthcheckRoute godoc
// @Summary Healthcheck
// @Description Simple healthcheck
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200
// @Header 200 {string} Token "qwerty"
func addHealthceckRoute(rg *gin.RouterGroup) {
	ping := rg.Group("/ping")

	ping.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
}

// Defines /v1/cache endpoint
func addCacheRoute(rg *gin.RouterGroup) {
	caches := rg.Group("/cache")

	// Handler for GET /v1/cache/:domain that check if cache for domain exists and returns cache size
	// User middleware checkContentType, checkAPIKey and checkDomainName on this call
	caches.GET("/:domain", validate.CheckContentType(), validate.CheckAPIKey(APIKey), validate.CheckDomainName(), func(c *gin.Context) {
		// Validate request by calling is isRequestValid and passing context
		if isRequestValid(c) {
			// If request is valid, try to find cache folder
			cacheSize, err := checkCache(c.Param("domain"))
			if err != nil {
				// If we fail to get chache folder info, we are returning error that caused cache not beeing deleted
				if err.Error() == "cache_folder_not_found" {
					c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
				}
			} else {
				// If cache lookup was successfull, we return cache size and 200
				c.JSON(http.StatusOK, gin.H{"status": "success", "domain": c.Param("domain"), "cache_size": cacheSize})
			}

		}
	})

	// Handler for DELETE /v1/cache/:domain that deletes contents of cache folder
	// User middleware checkContentType, checkAPIKey and checkDomainName on this call
	caches.DELETE("/:domain", validate.CheckContentType(), validate.CheckAPIKey(APIKey), validate.CheckDomainName(), func(c *gin.Context) {
		if isRequestValid(c) {
			// If request is valid, try to delete cache
			err := deleteCache(c.Param("domain"))
			if err != nil {
				// If we fail to delete cache, we are returning error that caused cache not beeing deleted
				if err.Error() == "cache_folder_not_found" {
					c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
				}
			} else {
				// If deletion was successfull, we return 200
				c.JSON(http.StatusOK, gin.H{"status": "success", "domain": c.Param("domain"), "msg": "cache_purged"})
			}

		}
	})
}
