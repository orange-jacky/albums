package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/common/util"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"net/http"
)

func AlbumManage(c *gin.Context) {
	user := util.GetUserName(c)
	album := util.GetAlbumName(c)
	mongo_album := util.GetAlbum()

	begin := GetMills()

	//
	resp := &data.Response{}
	action := c.Param("action")

	switch action {
	case "insert":
		var str string
		err := mongo_album.HasAblum(user, album)
		if err == nil {
			str = fmt.Sprintf("user had %v", user, album)
		} else {
			err := mongo_album.Insert(user, album)
			if err == nil {
				str = fmt.Sprintf("%v create %v success", user, album)
			} else {
				resp.Status = -1
				str = fmt.Sprintf("%v create %v fail, %v", user, album, err)
			}
		}
		resp.StatusDescription = str
	case "delete":
		var str string
		err := mongo_album.Delete(user, album)
		if err == nil {
			str = fmt.Sprintf("%v delete %v success", user, album)
			//删除完相册后,需要把相册里image信息删掉
			imageinfo := util.GetImageInfo()
			err := imageinfo.DeleteByUserAlbum(user, album)
			if err != nil {
				str = fmt.Sprintf("%v, delete imageinfo fail,%v", str, err)
			} else {
				str = fmt.Sprintf("%v, delete imageinfo success", str)
			}
		} else {
			str = fmt.Sprintf("%v delete %v fail, %v", user, album, err)
			resp.Status = -2
		}
		resp.StatusDescription = str
	case "get":
		rets, err := mongo_album.GetAlbums(user)
		if err == nil {
			resp.Data = rets
			resp.Total = len(rets)
			str := fmt.Sprintf("%v get albums success", user)
			resp.StatusDescription = str
		} else {
			str := fmt.Sprintf("%v get album fail, %v", user, err)
			resp.StatusDescription = str
			resp.Status = -3
		}
	}
	resp.Cost = GetMills() - begin
	if resp.Data == nil {
		resp.Data = make([]string, 0)
	}
	c.JSON(http.StatusOK, resp)
}
