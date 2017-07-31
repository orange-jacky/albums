package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/router"
	"github.com/orange-jacky/albums/util"
	"log"
	"net/http"
	"os"
	"time"
)

func usage(programName string) {
	fmt.Println(`
usage: 
	albums [configure file]

eg: albums conf/conf.xml
		`)
}

func main() {
	if len(os.Args) != 2 {
		usage(os.Args[0])
		os.Exit(-1)
	}
	//加载配置文件
	configure := util.Configure(os.Args[1])
	//启动日志单实例
	mylog := util.Mylog(configure.Log.File)
	defer mylog.Flush()

	//创建jobqueue
	jobqueue := util.JobQueue()
	jobqueue.Start()
	defer jobqueue.Stop()

	authMiddleware := GetAuthMiddleware()

	//设置gin模式
	gin.SetMode(gin.ReleaseMode)
	//配置路由
	r := gin.Default()
	r.POST("/signup", router.SignUp)
	r.POST("/login", authMiddleware.LoginHandler)
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.POST("/upload", router.UpLoad)
		auth.POST("/download", router.DownLoad)
		auth.POST("/search", router.Search)
	}
	//起一个http服务器
	server := fmt.Sprintf("%s:%s", configure.Gin.Host, configure.Gin.Port)
	s := &http.Server{
		Addr:         server,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatalln(s.ListenAndServe())
}
