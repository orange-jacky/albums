package db

import (
	"fmt"
	. "github.com/orange-jacky/albums/data"
	"github.com/orange-jacky/albums/util"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"os"
)

type MongoGridfs struct {
	Mongo
	GridFS *mgo.GridFS
}

func NewMongoGridfs() *MongoGridfs {
	return &MongoGridfs{}
}

func (m *MongoGridfs) OpenTable(table string) error {
	m.GridFS = m.Db.GridFS("fs")
	return nil
}

func (m *MongoGridfs) Query(query interface{}) (images []map[string]interface{}, err error) {
	record := make(map[string]interface{})
	iter := m.GridFS.Find(query).Iter()
	for iter.Next(&record) {
		//for debug
		//fmt.Println("record", record)
		images = append(images, record)
		record = make(map[string]interface{})
	}
	if err = iter.Close(); err != nil {
		return images, err
	}
	return images, nil
}

func (m *MongoGridfs) Insert(dir string, docs Images) error {
	for _, image := range docs {
		//文件路径
		filename := fmt.Sprintf("%s%s%s", dir, util.DirSeg(), image.Filename)
		//
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
		gridfile.SetId(image.Id)
		gridfile.SetName(image.Filename)
		gridfile.SetContentType(image.ContentType)
		gridfile.SetMeta(image.Metadata)

		if _, err := gridfile.Write(content); err != nil {
			return err
		}
	}
	return nil
}
