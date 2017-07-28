package util

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// Log 保存日志配置信息
type Log struct {
	File string `xml:"file"`
}

// Gin gin配置
type Gin struct {
	Host string `xml:"host"`
	Port string `xml:"port"`
}

// Feature 实时提取特征向量配置
type MyFeature struct {
	Host string `xml:"host"`
	Port string `xml:"port"`
}

// Nginx nginx-gridfs 获取图片配置
type Nginx struct {
	Host   string `xml:"host"`
	Port   string `xml:"port"`
	Router string `xml:"router"`
}

// Mongo mongo图片库和特征库配置
type MyMongo struct {
	Hosts string `xml:"hosts"`
	//Image 图片库
	Image struct {
		Db string `xml:"db"` //图片库名称
	} `xml:"image"`
	// Feature 特征库
	Feature struct {
		Db         string `xml:"db"`         //特征库名称
		Collection string `xml:"collection"` //特征表
	} `xml:"feature"`
}

type configure struct {
	XMLName xml.Name  `xml:"configure"`
	Log     Log       `xml:"log"`
	Gin     Gin       `xml:"gin"`
	Feature MyFeature `xml:"feature"`
	Nginx   Nginx     `xml:"nginx"`
	Mongo   MyMongo   `xml:"mongo"`
}

var (
	conf      *configure
	conf_once sync.Once
)

//Configure 载入xml配置文件
func Configure(file string) *configure {
	conf_once.Do(func() {
		conf = &configure{}
		if err := conf.init(file); err != nil {
			log.Fatalln(err)
		}
	})
	return conf
}

//init 载入xml配置文件
func (c *configure) init(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("error open file %s fail,%v", file, err)
	}
	defer fd.Close()
	content, _ := ioutil.ReadAll(fd)
	xml.Unmarshal(content, c)
	return nil
}
