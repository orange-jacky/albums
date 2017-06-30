package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// DownLoad 下载相册或者从相册里下载一张或者多张照片
func DownLoad(c *gin.Context) {
	c.String(http.StatusOK, "%s", "download")
}
