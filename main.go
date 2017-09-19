package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/router"
	"github.com/orange-jacky/albums/util"
)

var stoper []data.Stoper

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
	Init()
	defer Release()

	//设置gin模式
	gin.SetMode(gin.ReleaseMode)
	//配置路由
	r := gin.Default()
	r.POST("/signup", router.SignUp)
	authMiddleware := router.GetAuthMiddleware()
	r.POST("/login", authMiddleware.LoginHandler)
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/test", router.SignIn)
		auth.POST("/upload", router.UpLoad)
		auth.POST("/download", router.DownLoad)
		auth.POST("/search", router.Search)
		auth.POST("/managealbum/:action", router.AlbumManage)
		auth.POST("/delete", router.Delete)
		auth.POST("/deeplearning", router.DeepLearning)
		auth.POST("/objectdetection_dl", router.ObjectDetectionDL)
	}

	server := fmt.Sprintf("%s:%s", configure.Gin.Host, configure.Gin.Port)
	//起一个http服务器
	s := &http.Server{
		Addr:         server,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	go func(s *http.Server) {
		log.Printf("[Main] http server start\n")
		err := s.ListenAndServe()
		log.Printf("[Main] http server stop (%+v)\n", err)
	}(s)
	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	for {
		select {
		case sig := <-signals:
			log.Println("[Main] Catch signal", sig)
			//平滑关闭server
			err := s.Shutdown(context.Background())
			log.Printf("[Main] start gracefully shuts down http serve %+v", err)
			return
		}
	}
}

func Init() {
	//启动日志单实例
	mylog := util.Mylog()
	stoper = append(stoper, mylog)
	//创建jobqueue
	jobqueue := util.JobQueue()
	jobqueue.Start()
	stoper = append(stoper, jobqueue)
	//user
	user := util.MongoUser()
	stoper = append(stoper, user)
	//album
	album := util.MongoAlbum()
	stoper = append(stoper, album)
	//image
	image := util.MongoImage()
	stoper = append(stoper, image)
	//imageinfo
	imageInfo := util.MongoImageInfo()
	stoper = append(stoper, imageInfo)
	//feature
	f := util.Service_feature()
	stoper = append(stoper, f)
}

func Release() {
	for _, v := range stoper {
		v.Stop()
	}
}
