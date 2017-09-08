package router

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/common/jobqueue"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// UpLoad 上传图片到相册
func UpLoad(c *gin.Context) {
	begin := util.GetMills()

	user := util.GetUserName(c)
	album := util.GetAlbumName(c)
	//先缓存文件
	images, imageinfos, dir, err := cacheFile(user, album, c)
	if err != nil {
		clearCache(dir)
		mylog := util.GetMylog()
		mylog.Errorf("upload cache images fail, %v", err)
		resp := data.Response{Status: -1, Data: fmt.Sprintf("upload cache images fail, %v", err)}
		c.JSON(http.StatusOK, resp)
		return
	}
	//图片入库
	err = handleImage(images)
	if err != nil {
		clearCache(dir)
		mylog := util.GetMylog()
		mylog.Errorf("upload image insert mongo fail, %v", err)
		resp := data.Response{Status: -2, Data: fmt.Sprintf("upload image insert mongo fail, %v", err)}
		c.JSON(http.StatusOK, resp)
		return
	}
	//提取特征和入库服务
	handleImageInfos(dir, imageinfos)
	util.HandleUrl(imageinfos)

	resp := data.Response{}
	resp.Data = imageinfos
	resp.Total = len(imageinfos)
	resp.Cost = util.GetMills() - begin

	c.JSON(http.StatusOK, resp)
	//c.String(http.StatusOK, "%s", "upload")
}

// cacheFile 先把上传文件缓存到本地磁盘
func cacheFile(user, album string, c *gin.Context) (images data.Images, imageinfos data.ImageInfos,
	dir string, err error) {
	//根据上传时间,生成上传的唯一目录
	dir = util.GetDir(user, album, fmt.Sprintf("%d", util.GetNano()))

	r := c.Request
	//POST takes the uploaded file(s) and saves it to disk.
	//parse the multipart form in the request
	err = r.ParseMultipartForm(100000)
	if err != nil {
		return images, imageinfos, dir, err
	}
	//get a ref to the parsed multipart form
	m := r.MultipartForm
	//get the *fileheaders
	files := m.File["images"] //表单的name,id
	//创建临时缓存目录
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return images, imageinfos, dir, fmt.Errorf("create dir %s fail,%+v", dir, err)
	}
	//取所有上传文件
	for _, v := range files {
		//for each fileheader, get a handle to the actual file
		src, err := v.Open()
		if err != nil {
			continue
		}
		defer src.Close()
		filename := v.Filename
		path := filepath.Join(dir, filename)
		//create destination file making sure the path is writeable.
		dst, err := os.Create(path)
		if err != nil {
			continue
		}
		defer dst.Close()
		io.Copy(dst, src)

		//计算文件md5值
		content, err := ioutil.ReadFile(path)

		if err != nil || len(content) == 0 {
			mylog := util.GetMylog()
			mylog.Errorf("file:%v path:%v content:%v, err:%v", filename, path, len(content), err)
			continue
		}
		v_md5 := fmt.Sprintf("%x", md5.Sum(content))
		//image
		image := &data.Image{Filepath: path, Md5: v_md5}
		images = append(images, image)
		//imageinfo
		info := &data.ImageInfo{}
		info.User = user
		info.Album = album
		info.Filename = filename
		info.Filepath = path
		info.Updatetime = util.GetMills()
		info.Md5 = v_md5
		info.Url = v_md5
		imageinfos = append(imageinfos, info)
	}
	return images, imageinfos, dir, nil
}

//清理掉缓存文件
func clearCache(dir string) {
	os.RemoveAll(dir)
}

//图片入库
func handleImage(images data.Images) error {
	image := util.GetImage()
	err := image.Insert(images)
	return err
}

//发送给提取特征和imageinfo入库服务
type Input struct {
	dir        string
	imageinfos data.ImageInfos
}

func handleImageInfos(dir string, imageinfos data.ImageInfos) {
	input := &Input{dir, imageinfos}
	jobqueue := util.GetJobQueue()
	job := Job{Input: input, Handler: saveImageInfos}
	jobqueue.Push(job)
}

func saveImageInfos(input, output interface{}) {
	if input, ok := input.(*Input); ok {
		dir := input.dir
		defer clearCache(dir)
		service := util.GetService_feature()
		for _, v := range input.imageinfos {
			content, err := ioutil.ReadFile(v.Filepath)
			if err != nil {
				continue
			}
			features := service.Extract(content)
			v.Features = features
		}
		save := util.GetImageInfo()
		err := save.Insert(input.imageinfos)
		if err != nil {
			mylog := util.GetMylog()
			mylog.Infof("saveImageInfos fail,%v", err)
		}
	}
}
