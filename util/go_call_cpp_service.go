package util

import (
	"fmt"
	"github.com/orange-jacky/albums/imagehandle"
	"log"
	"sync"
)

//包装了一下go写的thrit客户端,所有go调用c++服务都从这个go服务发起
type service_feature struct {
	Hosts string
}

var (
	s_f      *service_feature
	s_f_once sync.Once
)

func Service_feature() *service_feature {
	s_f_once.Do(func() {
		s_f = &service_feature{}
		if err := s_f.Init(); err != nil {
			log.Fatalln(err)
		}
	})
	return s_f
}

//对外方法,使用时,先init,再start,退出时stop
func (m *service_feature) Init() error {
	conf := GetConfigure()
	m.Hosts = fmt.Sprintf("%s:%s", conf.Feature.Host, conf.Feature.Port)
	return nil
}

//提取图片hsv特征
func (m *service_feature) Extract(image []byte) (features []float64) {
	f, _ := imagehandle.GetImgFeature(m.Hosts, image)
	return f
}

//深度学习
func (m *service_feature) DeepLearning(image []byte) (r *imagehandle.Result_) {
	f, _ := imagehandle.DeepLearning(m.Hosts, image)
	return f
}

//深度学习做物体检测
func (m *service_feature) ObjectDetectionDL(image []byte) (r *imagehandle.Result_) {
	f, _ := imagehandle.ObjectDetectionDL(m.Hosts, image)
	return f
}

func (m *service_feature) Stop() {
}

func GetService_feature() *service_feature {
	return s_f
}
