package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/data"
	"path/filepath"
)

func GetUserName(c *gin.Context) string {
	user := c.PostForm("username")
	if user == "" {
		user = "default"
	}
	return user
}

func GetPassword(c *gin.Context) string {
	return c.PostForm("password")
}

func GetAlbumName(c *gin.Context) string {
	album := c.PostForm("album")
	if album == "" {
		album = "default"
	}
	return album
}

func HandleUrl(imageinfos data.ImageInfos) {
	url := GetUrl()
	for _, info := range imageinfos {
		info.Url = fmt.Sprintf("%s/%v", url, info.Md5)
	}
}

func GetDir(user, album, time string) string {
	return filepath.Join("filecache", user, fmt.Sprintf("%v-%v", album, time))
}

//生成访问图片路由
func GetUrl() string {
	conf := GetConfigure()
	return fmt.Sprintf("%s:%s/%s", conf.Nginx.HostInter, conf.Nginx.Port, conf.Nginx.Router)
}
