package util

import (
	"fmt"
	. "github.com/orange-jacky/albums/data"
	. "github.com/orange-jacky/albums/db"
	"log"
	"sync"
)

type mongoImage struct {
	gridfs *MongoGridfs
}

var (
	image      *mongoImage
	image_once sync.Once
)

func MongoImage() *mongoImage {
	image_once.Do(func() {
		image = &mongoImage{}
		if err := image.Init(); err != nil {
			log.Fatalln(err)
		}
	})
	return image
}

//对外方法,使用时,先init,再start,退出时stop
func (m *mongoImage) Init() error {
	conf := GetConfigure()

	mongo := MongoClient{}
	mongo.Hosts = conf.Mongo.Hosts
	mongo.Database = conf.Mongo.Image.Db
	mongo.Collection = conf.Mongo.Image.Collection
	gridfs := &MongoGridfs{MongoClient: mongo}

	m.gridfs = gridfs

	return nil
}

func (m *mongoImage) Insert(images Images) error {
	if err := m.gridfs.Connect(); err != nil {
		return fmt.Errorf("connect %v", err)
	}
	defer m.gridfs.Close()
	m.gridfs.DB()
	m.gridfs.C()
	return m.gridfs.Insert(images)
}

func (m *mongoImage) Stop() {

}

func GetImage() *mongoImage {
	return image
}
