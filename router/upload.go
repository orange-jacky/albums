package router

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/db"
	"github.com/orange-jacky/albums/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UpLoad 上传图片到相册
func UpLoad(c *gin.Context) {
	conf := util.Configure("")
	//log := util.Mylog("")

	//新建一个图片库连接
	gridfs := db.NewMongoGridfs()
	gridfs.Connect(conf.Mongo.Hosts, conf.Mongo.Image.Db)
	gridfs.OpenDb(conf.Mongo.Image.Db)
	gridfs.OpenTable("fs")
	defer gridfs.Close()

	//先缓存文件
	images, _ := cacheFile(c)
	defer clearCache(c)
	//提取图片信息
	getImageInfo(c, &images)
	//入库
	gridfs.Insert(images)

	resp := data.Response{}
	resp.Data = images
	c.JSON(http.StatusOK, resp)
	//c.String(http.StatusOK, "%s", "upload")
}

func getUser(c *gin.Context) string {
	user := c.PostForm("user")
	if user == "" {
		user = "default"
	}
	return user
}

func getAlbum(c *gin.Context) string {
	album := c.PostForm("album")
	if album == "" {
		album = "default"
	}
	return album
}

func clearCache(c *gin.Context) {
	user := getUser(c)
	album := getAlbum(c)
	os.RemoveAll(fmt.Sprintf("%s%s%s", user, util.DirSeg(), album))
}

// cacheFile 先把上传文件缓存到本地磁盘
func cacheFile(c *gin.Context) (images data.Images, err error) {
	r := c.Request
	//POST takes the uploaded file(s) and saves it to disk.
	//parse the multipart form in the request
	err = r.ParseMultipartForm(100000)
	if err != nil {
		return images, err
	}
	//get a ref to the parsed multipart form
	m := r.MultipartForm
	//get the *fileheaders
	files := m.File["images"] //表单的name,id

	//post 没有文件直接返回
	if len(files) == 0 {
		return images, nil
	}
	//创建临时缓存目录
	user := getUser(c)
	album := getAlbum(c)
	dir := fmt.Sprintf("%s%s%s", user, util.DirSeg(), album)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return images, err
	}
	//取所有上传文件
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		src, err := files[i].Open()
		if err != nil {
			continue
		}
		defer src.Close()
		//create destination file making sure the path is writeable.
		dst, err := os.Create(fmt.Sprintf("%s%s%s", dir, util.DirSeg(), files[i].Filename))
		if err != nil {
			continue
		}
		defer dst.Close()
		io.Copy(dst, src)
	}
	return images, nil
}

func getImageInfo(c *gin.Context, images *data.Images) {
	user := getUser(c)
	album := getAlbum(c)
	dir := fmt.Sprintf("%s%s%s", user, util.DirSeg(), album)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//跳过目录
		if info.IsDir() {
			return nil
		}
		filename := filepath.Base(path)
		img := &data.Imagedata{}
		img.User = user
		img.Album = album
		img.Name = filename
		sli := strings.Split(filename, ".")
		ext := strings.ToLower(sli[len(sli)-1])
		img.Type = ext
		img.Updatetime = util.GetMills()
		img.Filename = filename
		img.ContentType = fmt.Sprintf("image/%s", ext)
		img.Id = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s%s%v", user, album, filename, img.Updatetime))))

		*images = append(*images, img)
		return nil
	})
}
