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
	GridFS  *mgo.GridFS
	Results Images
}

func NewMongoGridfs() *MongoGridfs {
	return &MongoGridfs{}
}

func (m *MongoGridfs) OpenTable(table string) error {
	m.GridFS = m.Db.GridFS("fs")
	return nil
}

func (m *MongoGridfs) Query(query interface{}) error {
	record := new(Imagedata)
	iter := m.GridFS.Find(query).Iter()
	for iter.Next(record) {
		m.Results = append(m.Results, record)
		*record = Imagedata{} //重置变量
	}
	if err := iter.Close(); err != nil {
		return err
	}
	return nil
}

func (m *MongoGridfs) Insert(docs ...interface{}) error {
	for _, doc := range docs {
		//类型诊断
		if images, ok := doc.(Images); ok {
			for _, image := range images {
				//文件路径
				filename := fmt.Sprintf("%s%s%s%s%s", image.User, util.DirSeg(), image.Album, util.DirSeg(), image.Filename)
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
		} else {
			fmt.Println("type assert fail")
		}
	}
	return nil
}