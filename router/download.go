package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"net/http"
)

// DownLoad 下载相册或者从相册里下载一张或者多张照片
func DownLoad(c *gin.Context) {
	begin := util.GetMills()

	user := util.GetUserName(c)
	album := util.GetAlbumName(c)

	imageinfo := util.GetImageInfo()
	results, err := imageinfo.GetImageInfos(user, album)

	//返回数据
	resp := data.Response{}
	if err != nil {
		resp.Status = -1
		resp.Data = fmt.Sprintf("%v", err)
	} else {
		util.HandleUrl(results)
		resp.Data = results
		resp.Total = len(results)
	}

	resp.Cost = util.GetMills() - begin

	c.JSON(http.StatusOK, resp)
	//c.String(http.StatusOK, "%s", "download")
}
