package router

import (
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/common/util"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
)

// Search 以图搜图
func Search(c *gin.Context) {
	begin := GetMills()

	resp := data.Response{}
	//获取图片内容
	image_content, err := getsSearchFile(c)
	if err != nil {
		resp.Data = err
		c.JSON(http.StatusOK, resp)
		return
	}
	//提取特征
	s := util.GetService_feature()
	vect := s.Extract(image_content)

	//查特征库,找到对比数据
	imageinfos, err := queryImageInfo(c)
	if err != nil {
		resp.Data = err
		c.JSON(http.StatusOK, resp)
		return
	}
	//做卡方相似计算
	ret := histogram(vect, imageinfos)

	//分页
	page := util.GetPage(c)
	size := util.GetPageSize(c)
	start := page * size
	end := start + size

	ret = ret[start:end]
	//
	util.HandleUrl(ret)
	resp.Data = ret
	resp.Total = len(ret)
	resp.Cost = GetMills() - begin

	c.JSON(http.StatusOK, resp)
	//c.String(http.StatusOK, "search")
}

// getsSearchFile 获取搜索图片内容
func getsSearchFile(c *gin.Context) (image []byte, err error) {
	r := c.Request
	//POST takes the uploaded file(s) and saves it to disk.
	//parse the multipart form in the request
	err = r.ParseMultipartForm(100000)
	if err != nil {
		return image, err
	}
	//get a ref to the parsed multipart form
	m := r.MultipartForm
	//get the *fileheaders
	files := m.File["image"] //表单的name,id

	//post 没有文件直接返回
	if len(files) == 0 {
		return image, nil
	}
	//取第一张图片内容
	//for each fileheader, get a handle to the actual file
	src, err := files[0].Open()
	if err != nil {
		return image, err
	}
	defer src.Close()
	return ioutil.ReadAll(src)
}

//查imageinfo表
func queryImageInfo(c *gin.Context) (imageInfos data.ImageInfos, err error) {
	user := util.GetUserName(c)
	album := util.GetAlbumName(c)
	sort := []string{"updatetime"}
	skip := 0
	limit := 0

	s := util.GetImageInfo()
	return s.GetImageInfos(user, album, sort, skip, limit)
}

func histogram(search_vector []float64, imageinfos data.ImageInfos) (ret data.ImageInfos) {
	type A struct {
		distance  float64
		imageinfo *data.ImageInfo
	}
	var sli []*A
	for _, info := range imageinfos {
		d := chi2_distance(search_vector, info.Features)
		a := &A{d, info}
		sli = append(sli, a)
	}
	sort.Slice(sli, func(i, j int) bool { return sli[i].distance < sli[j].distance })

	for _, a := range sli {
		ret = append(ret, a.imageinfo)
		//fmt.Println(a.distance, a.imageinfo)
	}
	return ret
}

// chi-square  卡方检验
func chi2_distance(a, b []float64) (d float64) {
	//fmt.Println("a:", a)
	//fmt.Println("b:", b)
	c := zip(a, b)
	//fmt.Println("c:", c)
	for _, ab := range c {
		a, b := ab[0], ab[1]
		//	fmt.Println("a=", a, "b=", b)
		v := math.Pow(a-b, 2) / (a + b + 1e-10)
		//	fmt.Println("v=", v)
		d += v
	}
	d *= 0.5
	return d
}

func zip(a, b []float64) (c [][]float64) {
	len_a := len(a)
	len_b := len(b)
	min := int(math.Min(float64(len_a), float64(len_b)))
	for i := 0; i < min; i++ {
		sli := make([]float64, 0)
		sli = append(sli, a[i], b[i])
		//fmt.Printf("a[%d]=%v, b[%d]=%v, sli=%v\n", i, a[i], i, b[i], sli)
		c = append(c, sli)
	}
	return c
}
