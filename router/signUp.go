package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/common/util"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"net/http"
)

func SignUp(c *gin.Context) {
	begin := GetMills()

	user := util.GetUserName(c)
	album := util.GetAlbumName(c)
	passwd := util.GetPassword(c)

	u := util.GetUser()
	status, description := u.CheckUser(user, passwd)
	if status == data.NEW_USER { //新用户创建默认相册
		mongo_album := util.GetAlbum()
		err := mongo_album.Insert(user, album)
		if err != nil {
			description = fmt.Sprintf("%v,%v", description, err)
		}
	}
	resp := data.Response{Status: status, Data: description}
	resp.Cost = GetMills() - begin

	c.JSON(http.StatusOK, resp)
}
