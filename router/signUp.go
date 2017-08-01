package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/db"
	"github.com/orange-jacky/albums/util"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func SignUp(c *gin.Context) {
	r := c.Request

	conf := util.Configure("")
	mongoDB := db.NewMongo()
	err := mongoDB.Connect(conf.Mongo.Hosts, conf.Mongo.User.Db)
	if err != nil {
		resp := data.Response{Status: -2, Data: "连接mongo失败"}
		c.JSON(http.StatusOK, resp)
		return
	}
	mongoDB.OpenDb(conf.Mongo.User.Db)
	mongoDB.OpenTable(conf.Mongo.User.Collection)
	defer mongoDB.Close()

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if len(username) == 0 || len(password) == 0 {
		resp := data.Response{Status: -2, Data: "没有post username和password"}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		if mongoDB.FindUserOne(username) {
			resp := data.Response{Status: -2, Data: "用户已经注册过"}
			c.JSON(http.StatusOK, resp)
			return
		}
		mongoDB.Insert(bson.M{"username": username, "password": password})
		resp := data.Response{Status: 0, Data: fmt.Sprintf("成功注册用户 %s", username)}
		c.JSON(http.StatusOK, resp)
		return
	}
}
