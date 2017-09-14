package util

import (
	"fmt"
	. "github.com/orange-jacky/albums/data"
	. "github.com/orange-jacky/albums/db"
	"gopkg.in/mgo.v2"
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

	collection := m.Mongo.GetCollection()
	for _, info := range imageInfos {
		q := bson.M{"user": info.User, "album": info.Album, "md5": info.Md5}
		_, err := collection.Upsert(q, info)
		if err != nil { //不存在
			return fmt.Errorf("insert %v fail,%v", info, err)
		}
	}
	return nil
}

func (m *mongoImageInfo) Delete(imageInfos ImageInfos) error {
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	for _, info := range imageInfos {
		q := bson.M{"user": info.User, "album": info.Album, "md5": info.Md5}
		err := collection.Remove(q)
		if err != nil && err != mgo.ErrNotFound {
			return fmt.Errorf("delete %v  fail,%v", q, err)
		}
	}
	return nil
}

func (m *mongoImageInfo) DeleteByUserAlbum(user, album string) error {
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()
	q := bson.M{"user": user, "album": album}
	_, err := collection.RemoveAll(q)
	if err != nil {
		return fmt.Errorf("%v delete %v  fail,%v", user, album, err)
	}
	return nil
}

func (m *mongoImageInfo) GetImageInfos(user, album string, sort []string,
	skip, limit int) (rets ImageInfos, err error) {

	if err := m.Mongo.Connect(); err != nil {
		return rets, fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	var q bson.M
	if album == "" {
		q = bson.M{"user": user}
	} else {
		q = bson.M{"user": user, "album": album}
	}
	query := collection.Find(q)
	if sort != nil && len(sort) > 0 {
		query = query.Sort(sort...)
	}
	query = query.Skip(skip)
	if limit > 0 {
		query = query.Limit(limit)
	}
	iter := query.Iter()
	if err := iter.All(&rets); err != nil {
		return rets, fmt.Errorf("query %v %v sort:%v skip:%v limit:%v fail,%v", user, album,
			sort, skip, limit, err)
	}
	return rets, nil
}

func (m *mongoImageInfo) Stop() {

}

func GetImageInfo() *mongoImageInfo {
	return imageInfo
}
