package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
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

func handlerUrl(i interface{}) {
	url := util.GetUrl()
	switch i.(type) {
	case data.Images:
		for _, image := range i.(data.Images) {
			image.Filename = fmt.Sprintf("%s/%v", url, image.Filename)
		}
	case data.Features:
		for _, feature := range i.(data.Features) {
			feature.Filename = fmt.Sprintf("%s/%v", url, feature.Filename)
		}
	}
}
