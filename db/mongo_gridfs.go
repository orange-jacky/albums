package db

import (
	. "github.com/orange-jacky/albums/data"
	"gopkg.in/mgo.v2"
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
		//图片的md5相同,用新的覆盖旧的
		m.GridFS.RemoveId(image.Md5)

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
		gridfile.SetId(image.Md5)
		if _, err := gridfile.Write(content); err != nil {
			return err
		}
	}
	return nil
}
