package util

import (
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := c.Writer
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
		resp.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Encoding,Cache-Control,Content-Length,Accept-Encoding ")
	}
}
