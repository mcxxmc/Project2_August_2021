package webservice

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var blockedMethods = []string{"DELETE"}

func isStrInArray(s string, ss []string) bool {
	for _, e := range ss {
		if e == s {
			return true
		}
	}
	return false
}

// SetHeader sets the header to:
// 1. allow cross-origin (CORS)
func SetHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// Filter blocks unsupported methods like "DELETE".
func Filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if isStrInArray(method, blockedMethods) {
			fmt.Println("Unsupported HTTP Method.")
			c.Abort()
		} else {
			c.Next()
		}
	}
}
