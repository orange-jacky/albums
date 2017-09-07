package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"net/http"
)

func AlbumManage(c *gin.Context) {
	user := util.GetUserName(c)
	album := util.GetAlbumName(c)
	mongo_album := util.GetAlbum()

	//
	resp := &data.Response{}
	action := c.Param("action")

	switch action {
	case "insert":
		err := mongo_album.Insert(user, album)
		if err == nil {
			resp.Data = fmt.Sprintf("%v create %v success", user, album)
		} else {
			resp.Data = fmt.Sprintf("%v create %v fail, %v", user, album, err)
		}
	case "delete":
		err := mongo_album.Delete(user, album)
		if err == nil {
			resp.Data = fmt.Sprintf("%v delete %v success", user, album)
		} else {
			resp.Data = fmt.Sprintf("%v delete %v fail, %v", user, album, err)
		}
	case "get":
		rets, err := mongo_album.GetAlbums(user)
		if err == nil {
			resp.Data = rets
		} else {
			resp.Data = fmt.Sprintf("%v get album fail, %v", user, err)
		}
	}
	c.JSON(http.StatusOK, resp)
}
