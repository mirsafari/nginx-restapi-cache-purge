package validation

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

func CheckContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			c.Set("IsContentTypeValid", true)
		} else {
			c.Set("IsContentTypeValid", false)
		}

	}
}

func CheckAPIKey(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("APIKEY") == apiKey {
			c.Set("IsAuthorized", true)
		} else {
			c.Set("IsAuthorized", false)
		}
	}
}

func CheckDomainName() gin.HandlerFunc {
	return func(c *gin.Context) {
		RegExp := regexp.MustCompile(`^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|([a-zA-Z0-9][a-zA-Z0-9-_]{1,61}[a-zA-Z0-9]))\.([a-zA-Z]{2,6}|[a-zA-Z0-9-]{2,30}\.[a-zA-Z
			]{2,3})$`)
		if RegExp.MatchString(c.Param("domain")) {
			c.Set("IsDomainValid", true)
		} else {
			c.Set("IsDomainValid", false)
		}
	}
}
