package router

import (
	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) string {
	user := c.PostForm("user")
	if user == "" {
		user = "default"
	}
	return user
}

func getAlbum(c *gin.Context) string {
	album := c.PostForm("album")
	if album == "" {
		album = "default"
	}
	return album
}
