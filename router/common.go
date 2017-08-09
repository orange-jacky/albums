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
			image.Id = fmt.Sprintf("%s/%v", url, image.Id)
		}
	case data.Features:
		for _, feature := range i.(data.Features) {
			feature.Id = fmt.Sprintf("%s/%v", url, feature.Id)
		}
	case []map[string]interface{}:
		for _, m := range i.([]map[string]interface{}) {
			if v, ok := m["_id"]; ok {
				m["_id"] = fmt.Sprintf("%s/%s", url, v.(string))
			}
		}
	}
}
