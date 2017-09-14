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
func Delete(c *gin.Context) {
	begin := GetMills()

	user := util.GetUserName(c)
	album := util.GetAlbumName(c)
	md5_sli := util.GetDeleteMD5(c)

	var deletes data.ImageInfos
	for _, v := range md5_sli {
		info := &data.ImageInfo{User: user, Album: album, Md5: v}
		deletes = append(deletes, info)
	}
	imageinfo := util.GetImageInfo()
	err := imageinfo.Delete(deletes)

	//返回数据
	resp := data.Response{}
	if err != nil {
		resp.Status = -1
		resp.StatusDescription = fmt.Sprintf("%v", err)
	} else {
		resp.Data = deletes
		resp.Total = len(deletes)
		resp.StatusDescription = fmt.Sprintf("delete %v image successful", resp.Total)
	}

	resp.Cost = GetMills() - begin
	c.JSON(http.StatusOK, resp)
	//c.String(http.StatusOK, "%s", "download")
}
