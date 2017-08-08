package router

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/common"
	. "github.com/orange-jacky/albums/common/jobqueue"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/db"
	. "github.com/orange-jacky/albums/feature"
	"github.com/orange-jacky/albums/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// UpLoad 上传图片到相册
func UpLoad(c *gin.Context) {
	conf := util.Configure("")
	log := util.Mylog("")

	//新建一个图片库连接
	gridfs := db.NewMongoGridfs()
	gridfs.Connect(conf.Mongo.Hosts, conf.Mongo.Image.Db)
	gridfs.OpenDb(conf.Mongo.Image.Db)
	gridfs.OpenTable("fs")
	defer gridfs.Close()

	//先缓存文件
	images, dir, err := cacheFile(c)
	if err != nil {
		clearCache(dir)
		log.Errorf("upload cache images fail, %v", err)
		resp := data.Response{Status: -1, Data: fmt.Sprintf("upload cache images fail, %v", err)}
		c.JSON(http.StatusOK, resp)
	}
	//提取图片信息
	getImageInfo(c, getUser(c), getAlbum(c), dir, &images)
	//入库
	err = gridfs.Insert(dir, images)
	if err != nil {
		clearCache(dir)
		log.Errorf("upload image insert mongo fail, %v", err)
		resp := data.Response{Status: -2, Data: fmt.Sprintf("upload image insert mongo fail, %v", err)}
		c.JSON(http.StatusOK, resp)
	}
	//发送给提取特征和入库服务
	send2GetFeature(dir, images)

	resp := data.Response{}
	resp.Data = images
	c.JSON(http.StatusOK, resp)
	//c.String(http.StatusOK, "%s", "upload")
}

// cacheFile 先把上传文件缓存到本地磁盘
func cacheFile(c *gin.Context) (images data.Images, dir string, err error) {

	//根据上传时间,生成上传的唯一目录
	user := getUser(c)
	album := getAlbum(c)
	dir = util.GetDir(user, album, fmt.Sprintf("%d", GetNano()))

	r := c.Request
	//POST takes the uploaded file(s) and saves it to disk.
	//parse the multipart form in the request
	err = r.ParseMultipartForm(100000)
	if err != nil {
		return images, dir, err
	}
	//get a ref to the parsed multipart form
	m := r.MultipartForm
	//get the *fileheaders
	files := m.File["images"] //表单的name,id

	//post 没有文件直接返回
	if len(files) == 0 {
		return images, dir, fmt.Errorf("not post images")
	}
	//创建临时缓存目录
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return images, dir, fmt.Errorf("create dir %s fail,%+v", dir, err)
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
	return images, dir, nil
}

func getImageInfo(c *gin.Context, user, album, dir string, images *data.Images) {

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//处理walk时的error
		if err != nil {
			mylog := util.Mylog("")
			mylog.Errorf("filepath.walk %s %s fail,%s", dir, path, err.Error())
			return nil
		}
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

func clearCache(dir string) {
	os.RemoveAll(dir)
}

//发送给提取特征和入库服务
type A struct {
	dir    string
	images data.Images
}

func send2GetFeature(dir string, images data.Images) {
	a := A{dir, images}
	var feature []float64

	jobqueue := util.JobQueue()
	job := Job{Input: a, Output: feature, Handler: getImgFeature}
	jobqueue.Push(job)
}

func getImgFeature(input, output interface{}) {
	if a, ok := input.(A); ok {
		dir := a.dir
		defer clearCache(dir)

		var features data.Features
		//host
		conf := util.Configure("")
		hostport := fmt.Sprintf("%s:%s", conf.Feature.Host, conf.Feature.Port)

		//提取所有上传文件特征
		for _, image := range a.images {
			filename := fmt.Sprintf("%s%s%s", dir, util.DirSeg(), image.Metadata.Name)
			sli_b, err := ioutil.ReadFile(filename)
			if err != nil {
				continue
			}
			//提取特征
			var feature_vector []float64
			feature_vector, err = GetImgFeature(sli_b, hostport)

			feature := &data.Featuredata{}
			feature.Metadata = image.Metadata
			feature.Attr = image.Attr
			feature.Features = feature_vector

			features = append(features, feature)
		}
		//入库
		mongo := db.NewMongo()
		mongo.Connect(conf.Mongo.Hosts, conf.Mongo.Feature.Db)
		mongo.OpenDb(conf.Mongo.Feature.Db)
		mongo.OpenTable(conf.Mongo.Feature.Collection)
		defer mongo.Close()

		//log
		mylog := util.Mylog("")
		for _, v := range features {
			err := mongo.Insert(*v)
			if err != nil {
				mylog.Errorf("%s", err.Error())
			}
		}
	}
}
