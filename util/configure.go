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
type Feature struct {
	Host string `xml:"host"`
	Port string `xml:"port"`
}

// Nginx nginx-gridfs 获取图片配置
type Nginx struct {
	Host      string `xml:"host"`
	HostInter string `xml:"host_inter"`
	Port      string `xml:"port"`
	Router    string `xml:"router"`
}

type DBTable struct {
	Db         string `xml:"db"`
	Collection string `xml:"collection"`
}

// Mongo mongo图片库和特征库配置
type Mongo struct {
	Hosts string `xml:"hosts"`
	// 用户库
	User DBTable `xml:"user"`
	// 相册库
	Album DBTable `xml:"album"`
	// 图片库
	Image DBTable `xml:"image"`
	// 图片信息库
	ImageInfo DBTable `xml:"imageinfo"`
}

type configure struct {
	XMLName xml.Name `xml:"configure"`
	Log     Log      `xml:"log"`
	Gin     Gin      `xml:"gin"`
	Feature Feature  `xml:"feature"`
	Nginx   Nginx    `xml:"nginx"`
	Mongo   Mongo    `xml:"mongo"`
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

func GetConfigure() *configure {
	return conf
}
