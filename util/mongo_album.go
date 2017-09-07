package util

import (
	"fmt"
	. "github.com/orange-jacky/albums/data"
	. "github.com/orange-jacky/albums/db"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
)

type mongoAlbum struct {
	Mongo *MongoClient
}

var (
	album      *mongoAlbum
	album_once sync.Once
)

func MongoAlbum() *mongoAlbum {
	album_once.Do(func() {
		album = &mongoAlbum{}
		if err := album.Init(); err != nil {
			log.Fatalln(err)
		}
	})
	return album
}

//对外方法,使用时,先init,再start,退出时stop
func (m *mongoAlbum) Init() error {
	conf := GetConfigure()

	mongo := &MongoClient{}
	mongo.Hosts = conf.Mongo.Hosts
	mongo.Database = conf.Mongo.Album.Db
	mongo.Collection = conf.Mongo.Album.Collection
	m.Mongo = mongo
	return nil
}

func (m *mongoAlbum) Insert(user, album string) error {
	if user == "" || album == "" {
		return fmt.Errorf("input user or album empty")
	}
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	result := &Album{}
	q := bson.M{"user": user}
	c := m.Mongo.GetCollection()
	err := c.Find(q).One(result)
	if err != nil { //用户第一次创建相册
		a := &Album{}
		a.User = user
		a.Albums = append(a.Albums, album)
		if err := m.Mongo.Insert(a); err != nil {
			return fmt.Errorf("insert %v %v fail,%v", user, album, err)
		}
	} else { //用户已经创建过相册
		var find bool
		for _, v := range result.Albums {
			if v == album {
				find = true
				break
			}
		}
		if !find {
			selector := bson.M{"user": user}
			update := &Album{User: user}
			update.Albums = append(update.Albums, album)
			update.Albums = append(update.Albums, result.Albums...)
			if err := m.Mongo.Update(selector, update); err != nil {
				return fmt.Errorf("%v append %v %v", user, album, err)
			}
		}
	}
	return nil
}

func (m *mongoAlbum) Delete(user, album string) error {
	if user == "" || album == "" {
		return fmt.Errorf("input user or album empty")
	}
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	result := &Album{}
	q := bson.M{"user": user}
	c := m.Mongo.GetCollection()
	err := c.Find(q).One(result)
	if err != nil { //没找到用户
		return fmt.Errorf("user %v not exist", user)
	} else {
		selector := bson.M{"user": user}
		update := &Album{User: user}
		for _, v := range result.Albums {
			if v != album {
				update.Albums = append(update.Albums, v)
			}
		}
		if err := m.Mongo.Update(selector, update); err != nil {
			return fmt.Errorf("%v delete %v %v", user, album, err)
		}
	}
	return nil
}

func (m *mongoAlbum) GetAlbums(user string) (rets []string, err error) {
	if user == "" {
		return rets, fmt.Errorf("input user empty")
	}
	if err := m.Mongo.Connect(); err != nil {
		return rets, fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	result := &Album{}
	q := bson.M{"user": user}
	c := m.Mongo.GetCollection()
	err = c.Find(q).One(result)
	if err != nil {
		return rets, fmt.Errorf("query %v", err)
	}
	rets = append(rets, result.Albums...)
	return rets, nil
}

func (m *mongoAlbum) Stop() {

}

func GetAlbum() *mongoAlbum {
	return album
}
