package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

/* checkCache is a function that check if cache folder is present and returns its size
Args: domain (string) - domain name for which we check if cache exists
Returns: cacheSize (float32) - size of cache in MB
         err - error that gets returned to client
*/
func checkCache(domain string) (cacheSize float32, err error) {
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
func deleteCache(domain string) (err error) {
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
	for _, d := range cacheDir {
		os.RemoveAll(path.Join([]string{pathToCache, d.Name()}...))
	}

	return nil
}

/* exists returns whether the given directory exists
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
