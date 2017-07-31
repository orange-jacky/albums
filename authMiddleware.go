package main

import (
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"time"
	"github.com/orange-jacky/albums/db"
	"github.com/orange-jacky/albums/util"
	"gopkg.in/mgo.v2/bson"
	"github.com/orange-jacky/albums/data"
)


func GetAuthMiddleware() *jwt.GinJWTMiddleware{
	return &jwt.GinJWTMiddleware{
		Realm:      "test zone",
		Key:        []byte("secret key"),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		Authenticator: func(userName string, password string, c *gin.Context) (string, bool) {
			conf := util.Configure("")
			mongoDB := db.NewMongo()
			mongoDB.Connect(conf.Mongo.Hosts, conf.Mongo.User.Db)
			mongoDB.OpenTable("user")
			defer mongoDB.Close()
			if mongoDB.FindUserOne(bson.M{"userName":userName, "password":password}){
				return userName, true
			}

			return userName, false
		},
		Authorizator: func(userName string, c *gin.Context) bool {
			conf := util.Configure("")
			mongoDB := db.NewMongo()
			mongoDB.Connect(conf.Mongo.Hosts, conf.Mongo.User.Db)
			mongoDB.OpenTable("user")
			defer mongoDB.Close()
			if mongoDB.FindUserOne(bson.M{"userName":userName}) {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			resp := data.Response{Status:code, Data:message}
			c.JSON(code, resp)
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}
}
