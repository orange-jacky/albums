package db

import (
	. "github.com/orange-jacky/albums/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"os"
)

type MongoGridfs struct {
	MongoClient
	GridFS *mgo.GridFS
}

func (m *MongoGridfs) C() {
	fs := "fs"
	if m.Collection != "" {
		fs = m.Collection
	}
	m.GridFS = m.Session.Db.GridFS(fs)
}

func (m *MongoGridfs) Insert(images Images) error {
	for _, image := range images {
		//检查是否已经存在
		result := make(map[string]interface{})
		q := bson.M{"filename": image.Md5}
		if err := m.GridFS.Find(q).One(&result); err == nil { //找到md5相同的文件
			if id, ok := result["_id"]; ok {
				m.GridFS.RemoveId(id) //删除旧文件
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
		gridfile, err := m.GridFS.Create(filename)
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
