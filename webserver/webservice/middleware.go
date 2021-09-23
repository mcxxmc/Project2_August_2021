package webservice

import (
	"github.com/gin-gonic/gin"
)

// SetCORS sets the header to allow cross-origin (CORS)
func SetCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

// Filter blocks unsupported methods like "DELETE".
func Filter() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "POST, GET, UPDATE")
	}
}
