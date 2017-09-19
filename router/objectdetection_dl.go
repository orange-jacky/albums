package router

import (
	"github.com/gin-gonic/gin"
	. "github.com/orange-jacky/albums/common/util"
	"github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"net/http"
)

func ObjectDetectionDL(c *gin.Context) {
	begin := GetMills()

	resp := data.Response{}
	//获取图片内容
	image_content, err := getsSearchFile(c)
	if err != nil {
		resp.Status = -1
		resp.StatusDescription = err
		c.JSON(http.StatusOK, resp)
		return
	}
	//提取特征
	s := util.GetService_feature()
	vect := s.ObjectDetectionDL(image_content)

	resp.Data = vect
	resp.Total = 1
	resp.Cost = GetMills() - begin

	c.JSON(http.StatusOK, resp)
}
