package router

import (
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/util"
	"github.com/orange-jacky/albums/db"
	"log"
	"net/http"
	"github.com/orange-jacky/albums/data"
	"gopkg.in/mgo.v2/bson"
)

func SignUp(c *gin.Context) {
	defer func (){
		resp := data.Response{Status:-1, Data:"注册失败"}
		c.JSON(http.StatusServiceUnavailable, resp)
	}()
	r := c.Request

	conf := util.Configure("")

	mongoDB := db.NewMongo()
	mongoDB.Connect(conf.Mongo.Hosts, conf.Mongo.User.Db)
	mongoDB.OpenTable("user")
	defer mongoDB.Close()

	err := r.ParseMultipartForm(100000)
	if err != nil{
		errStr := "Parse form error"+err.Error()
		log.Fatal(errStr)
		panic(errStr)
	}

	userName := r.PostFormValue("userName")
	password := r.PostFormValue("password")

	if len(userName) == 0 && len(password) == 0{
		log.Fatal("Get userName or password error")
		panic("Get userName or password error")
	}else {
		if mongoDB.FindUserOne(bson.M{"userName":userName}){
			resp := data.Response{Status:-2, Data:"用户已经注册过"}
			c.JSON(http.StatusServiceUnavailable, resp)
		}
		mongoDB.Insert(bson.M{"userName":userName, "password": password})
		c.JSON(http.StatusOK, data.Response{})
	}
}
