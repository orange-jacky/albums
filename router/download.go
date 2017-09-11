package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/common/util"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"net/http"
)

// DownLoad 下载相册或者从相册里下载一张或者多张照片
func DownLoad(c *gin.Context) {
	begin := GetMills()

	user := util.GetUserName(c)
	album := util.GetAlbumName(c)

	page := util.GetPage(c)
	size := util.GetPageSize(c)
	sort := []string{"updatetime"}
	skip := page * size
	limit := size

	imageinfo := util.GetImageInfo()
	results, err := imageinfo.GetImageInfos(user, album, sort, skip, limit)

	//返回数据
	resp := data.Response{}
	if err != nil {
		resp.Status = -1
		resp.StatusDescription = fmt.Sprintf("%v", err)
	} else {
		util.HandleUrl(results)
		resp.Data = results
		resp.Total = len(results)
	}

	resp.Cost = GetMills() - begin

	c.JSON(http.StatusOK, resp)
	//c.String(http.StatusOK, "%s", "download")
}
