package util

import (
	"fmt"
	. "github.com/orange-jacky/albums/data"
	. "github.com/orange-jacky/albums/db"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
)

type mongoImageInfo struct {
	Mongo *MongoClient
}

var (
	imageInfo      *mongoImageInfo
	imageInfo_once sync.Once
)

func MongoImageInfo() *mongoImageInfo {
	imageInfo_once.Do(func() {
		imageInfo = &mongoImageInfo{}
		if err := imageInfo.Init(); err != nil {
			log.Fatalln(err)
		}
	})
	return imageInfo
}

//对外方法,使用时,先init,再start,退出时stop
func (m *mongoImageInfo) Init() error {
	conf := GetConfigure()

	mongo := &MongoClient{}
	mongo.Hosts = conf.Mongo.Hosts
	mongo.Database = conf.Mongo.ImageInfo.Db
	mongo.Collection = conf.Mongo.ImageInfo.Collection
	m.Mongo = mongo
	return nil
}

func (m *mongoImageInfo) Insert(imageInfos ImageInfos) error {
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	for _, info := range imageInfos {
		result := &ImageInfo{}
		q := bson.M{"user": info.User, "album": info.Album, "md5": info.Md5}
		c := m.Mongo.GetCollection()
		err := c.Find(q).One(result)
		if err != nil { //不存在
			if err := m.Mongo.Insert(info); err != nil {
				return fmt.Errorf("insert %v  fail,%v", info, err)
			}
		}
	}
	return nil
}

func (m *mongoImageInfo) GetImageInfos(user, album string, sort []string, skip, limit int) (rets ImageInfos, err error) {
	if user == "" {
		return rets, fmt.Errorf("input user empty")
	}
	if err := m.Mongo.Connect(); err != nil {
		return rets, fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	var q bson.M
	if album == "" {
		q = bson.M{"user": user}
	} else {
		q = bson.M{"user": user, "album": album}
	}
	query := m.Mongo.GetCollection().Find(q)
	if sort != nil && len(sort) > 0 {
		query = query.Sort(sort...)
	}
	query = query.Skip(skip)
	if limit > 0 {
		query = query.Limit(limit)
	}
	iter := query.Iter()
	if err := iter.All(&rets); err != nil {
		return rets, fmt.Errorf("query %v", err)
	}
	//if err := m.Mongo.Query(q, &rets); err != nil {
	//	return rets, fmt.Errorf("query %v", err)
	//}
	return rets, nil
}

func (m *mongoImageInfo) Stop() {

}

func GetImageInfo() *mongoImageInfo {
	return imageInfo
}
