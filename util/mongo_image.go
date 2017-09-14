package util

import (
	"fmt"
	. "github.com/orange-jacky/albums/data"
	. "github.com/orange-jacky/albums/db"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"os"
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

	gridfs := m.gridfs.GetGridFs()

	for _, image := range images {
		//检查是否已经存在
		result := make(map[string]interface{})
		q := bson.M{"filename": image.Md5}
		if err := gridfs.Find(q).One(&result); err == nil { //找到md5相同的文件
			if id, ok := result["_id"]; ok {
				gridfs.RemoveId(id) //删除旧文件
			}
		}
		filename := image.Filepath
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		//
		gridfile, err := gridfs.Create(filename)
		if err != nil {
			return err
		}
		defer gridfile.Close()
		gridfile.SetName(image.Md5)
		if _, err := gridfile.Write(content); err != nil {
			return err
		}
	}

	return nil
}

func (m *mongoImage) Delete(images Images) error {
	if err := m.gridfs.Connect(); err != nil {
		return fmt.Errorf("connect %v", err)
	}
	defer m.gridfs.Close()
	m.gridfs.DB()
	m.gridfs.C()

	gridfs := m.gridfs.GetGridFs()

	for _, image := range images {
		result := make(map[string]interface{})
		q := bson.M{"filename": image.Md5}
		if err := gridfs.Find(q).One(&result); err == nil { //找到md5相同的文件
			if id, ok := result["_id"]; ok {
				gridfs.RemoveId(id) //删除旧文件
			}
		}
	}
	return nil
}

func (m *mongoImage) Stop() {

}

func GetImage() *mongoImage {
	return image
}
