package router

import (
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/db"
	"github.com/orange-jacky/albums/util"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// DownLoad 下载相册或者从相册里下载一张或者多张照片
func DownLoad(c *gin.Context) {
	conf := util.Configure("")
	//log := util.Mylog("")

	//新建一个图片库连接
	gridfs := db.NewMongoGridfs()
	gridfs.Connect(conf.Mongo.Hosts, conf.Mongo.Image.Db)
	gridfs.OpenDb(conf.Mongo.Image.Db)
	gridfs.OpenTable("fs")
	defer gridfs.Close()

	//查数据库
	user := getUser(c)
	album := getAlbum(c)
	query := bson.M{"metadata.user": user, "metadata.album": album}
	images, _ := gridfs.Query(query)
	//返回数据
	resp := data.Response{}
	resp.Data = images
	c.JSON(http.StatusOK, resp)
	//c.String(http.StatusOK, "%s", "download")
}
