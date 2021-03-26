package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// CacheGET Handler for GET /v1/cache/:domain fetches info about requested domain
// @Summary Cache get route
// @Description "Route for checking if cache exists and returing value of it"
// @ID caches.GET
// @Accept  application/json
// @Produce  application/json
// @Param domain path string true "Domain"
// @Success 200 {object} models.CacheEndpointResponseGET
// @Failure 400 {object} models.CacheEndpointResponseGET
// @Failure 404 {object} models.CacheEndpointResponseGET
// @Failure 500 {object} models.CacheEndpointResponseGET
// @Router /cache/{domain} [get]
func CacheGET(c *gin.Context, cachePath string) {
	// Validate request by calling is isRequestValid and passing context
	if isRequestValid(c) {
		// If request is valid, try to find cache folder
		cacheSize, err := checkCache(cachePath, c.Param("domain"))
		if err != nil {
			// If we fail to get chache folder info, we are returning error that caused cache not beeing deleted
			if err.Error() == "cache_folder_not_found" {
				c.JSON(http.StatusNotFound, gin.H{"status": "error", "domain": c.Param("domain"), "msg": err.Error(), "cache_size": 0})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "domain": c.Param("domain"), "msg": err.Error(), "cache_size": 0})
			}
		} else {
			// If cache lookup was successfull, we return cache size and 200
			c.JSON(http.StatusOK, gin.H{"status": "success", "domain": c.Param("domain"), "cache_size": cacheSize})
		}

	}
}

// CacheDELETE Handler for DELETE /v1/cache/:domain deletes cache folder content for requested domain
// @Summary Cache delete route
// @Description "Route that tries to delete cahce folder for given domain"
// @ID caches.DELETE
// @Accept  application/json
// @Produce  application/json
// @Param domain path string true "Domain"
// @Success 200 {object} models.CacheEndpointResponseDELETE
// @Failure 400 {object} models.CacheEndpointResponseDELETE
// @Failure 404 {object} models.CacheEndpointResponseDELETE
// @Failure 500 {object} models.CacheEndpointResponseDELETE
// @Router /cache/{domain} [delete]
func CacheDelete(c *gin.Context, cachePath string) {
	if isRequestValid(c) {
		// If request is valid, try to delete cache
		err := deleteCache(cachePath, c.Param("domain"))
		if err != nil {
			// If we fail to delete cache, we are returning error that caused cache not beeing deleted
			if err.Error() == "cache_folder_not_found" {
				c.JSON(http.StatusNotFound, gin.H{"status": "error", "domain": c.Param("domain"), "msg": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "domain": c.Param("domain"), "msg": err.Error()})
			}
		} else {
			// If deletion was successfull, we return 200
			c.JSON(http.StatusOK, gin.H{"status": "success", "domain": c.Param("domain"), "msg": "cache_purged"})
		}

	}
}

/* checkCache is a function that check if cache folder is present and returns its size
Args: domain (string) - domain name for which we check if cache exists
Returns: cacheSize (float32) - size of cache in MB
         err - error that gets returned to client
*/
func checkCache(cachePath string, domain string) (cacheSize float32, err error) {
	// Create variable that contains full path to cache
	pathToCache := cachePath + domain

	// Call exists function
	cacheIsThere, err := exists(pathToCache)

	// And check if return value from function is false
	if !(cacheIsThere) {

		// If return value is false and there are not errors, this means that we did not find cache folder
		if err == nil {
			// So we need to setup proper error that will be returned to client
			err = errors.New("cache_folder_not_found")
		}
		// And we exit function returing the error and 0 since we have 2 return values
		return 0, err
	}

	// If cache folder exists, we need to get the size of the folder, so we create new variable
	var size int64
	// And create a function that will try to get size of the folder
	err = filepath.Walk(pathToCache, func(_ string, info os.FileInfo, err error) error {
		// If we encouter error, we exit the function and return an error
		if err != nil {
			return err
		}
		// Else we append to the total size
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	// At the end we calculate the size in MB and return any errors that might occur
	return float32(size) / 1024 / 1024, err
}

/* deleteCache is a function that deletes content of folder under given path
Args: domain (string) - domain name for which we delete cache
Returns: err - error that gets returned to client if deletion was unsucessfull
*/
func deleteCache(cachePath string, domain string) (err error) {
	// Create variable that contains full path to cache
	pathToCache := cachePath + domain

	// Call exists function
	cacheIsThere, err := exists(pathToCache)

	// And check if return value from function is false
	if !(cacheIsThere) {

		// If return value is false and there are not errors, this means that we did not find cache folder
		if err == nil {
			// So we need to setup proper error that will be returned to client
			err = errors.New("cache_folder_not_found")
		}
		// And we exit function returing the error
		return err
	}

	// If cache folder exists, we need delete all contents of the direcotry
	cacheDir, err := ioutil.ReadDir(pathToCache)
	if err != nil {
		return err
	}

	for _, d := range cacheDir {
		os.RemoveAll(path.Join([]string{pathToCache, d.Name()}...))
	}

	return nil
}

/* Exists returns whether the given directory exists
Args: path (string) - path to given file
Returns: bool - if path is found and it's a folder, we return true
         err  - if there has been an error in opening files, we return that error
*/
func exists(path string) (bool, error) {
	// Get info about the given path
	info, err := os.Stat(path)

	// Check if we have any errors opening location and we check if the location is direcory
	if err == nil && info.IsDir() {
		return true, nil
	} else if err == nil && !(info.IsDir()) {
		// If the location exists but it's not direcory, we need to respont to clinet
		err = errors.New("path_is_file_not_folder")
		return false, err
	}
	// We do additional checks with os.IsNotExist() function
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/* isRequestValid checks flags set by middleware if requests passes validation, if they fail we finish request and return proper error message
Args: *gin.Context (pointer to Context) - Context contaning all flags
Returns: bool - if request is valid, return value is true
*/
func isRequestValid(c *gin.Context) bool {
	if c.MustGet("IsContentTypeValid") != true {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "content_type_invalid"})
		return false
	}

	if c.MustGet("IsAuthorized") != true {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "msg": "api_key_invalid"})
		return false
	}

	if c.MustGet("IsDomainValid") != true {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": "domain_name_invalid"})
		return false
	}
	return true
}
