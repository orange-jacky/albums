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

	username := util.GetUserName(c)
	album := util.GetAlbumName(c)
	passwd := util.GetPassword(c)

	u := util.GetUser()
	a := util.GetAlbum()

	resp := &data.Response{}

	err := u.CheckUser(username)
	if err == nil {
		resp.Status = -1
		resp.StatusDescription = fmt.Sprintf("%v exist, please change a username", username)
	} else {
		if err == data.USER_NOT_EXIST {
			if err = u.NewUser(username, passwd); err == nil {
				if err = a.InsertDefault(username, album); err == nil {
					resp.StatusDescription = fmt.Sprintf("create a new user %v success", username)
				} else {
					resp.Status = -1
					resp.StatusDescription = fmt.Sprintf("%v", err)
				}
			} else {
				resp.Status = -1
				resp.StatusDescription = fmt.Sprintf("%v", err)
			}
		} else {
			resp.Status = -1
			resp.StatusDescription = fmt.Sprintf("%v", err)
		}
	}
	resp.Cost = GetMills() - begin
	c.JSON(http.StatusOK, resp)
}
