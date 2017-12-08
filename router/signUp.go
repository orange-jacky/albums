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

	count, err := u.CheckUser(username)
	if err != nil {
		resp.Status = -100
		resp.StatusDescription = fmt.Sprintf("%v", err)
	} else {
		if count == 0 {
			if err = u.NewUser(username, passwd); err == nil {
				if err = a.InsertDefault(username, album); err == nil {
					//插入test
					a.Insert(username, "test")
					resp.Status = 0
					resp.StatusDescription = fmt.Sprintf("create a new user %v success", username)

				} else {
					resp.Status = 1
					resp.StatusDescription = fmt.Sprintf("%v", err)
				}
			} else {
				resp.Status = -1
				resp.StatusDescription = fmt.Sprintf("%v", err)
			}
		} else if count == 1 {
			resp.Status = 2
			resp.StatusDescription = fmt.Sprintf("%v exists", username)
		} else {
			resp.Status = -2
			resp.StatusDescription = fmt.Sprintf("lots of %v exists ", username)
		}
	}
	resp.Cost = GetMills() - begin
	c.JSON(http.StatusOK, resp)
}
