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
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect db fail,%v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()
	collection := m.Mongo.GetCollection()

	result := &Album{}
	q := bson.M{"user": user}
	err := collection.Find(q).One(result)
	if err != nil {
		return fmt.Errorf("find %v from album fail, %v", user, err)
	} else {
		selector := q
		update := &Album{User: user}
		for _, v := range result.Albums {
			if v != album {
				update.Albums = append(update.Albums, v)
			}
		}
		if err := collection.Update(selector, update); err != nil {
			return fmt.Errorf("%v insert %v fail,%v", user, album, err)
		}
	}
	return nil
}

func (m *mongoAlbum) InsertDefault(user, album string) error {
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect db fail,%v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()
	a := &Album{}
	a.User = user
	a.Albums = append(a.Albums, album)
	if err := collection.Insert(a); err != nil {
		return fmt.Errorf("%v insert default %v fail,%v", user, album, err)
	}
	return nil
}

func (m *mongoAlbum) Delete(user, album string) error {
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect db fail,%v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	result := &Album{}
	q := bson.M{"user": user}
	err := collection.Find(q).One(result)
	if err != nil { //没找到用户
		return fmt.Errorf("find %v from album fail, %v", user, err)
	} else {
		selector := q
		update := &Album{User: user}
		for _, v := range result.Albums {
			if v != album {
				update.Albums = append(update.Albums, v)
			}
		}
		if err := collection.Update(selector, update); err != nil {
			return fmt.Errorf("%v delete %v %v", user, album, err)
		}
	}
	return nil
}

func (m *mongoAlbum) GetAlbums(user string) (rets []string, err error) {
	if err := m.Mongo.Connect(); err != nil {
		return rets, fmt.Errorf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	result := &Album{}
	q := bson.M{"user": user}
	err = collection.Find(q).One(result)
	if err != nil {
		return rets, fmt.Errorf("query %v %v", user, err)
	}
	rets = append(rets, result.Albums...)
	return rets, nil
}

func (m *mongoAlbum) HasAblum(user, album string) error {
	albums, err := m.GetAlbums(user)
	if err != nil {
		return err
	}
	for _, v := range albums {
		if v == album {
			return nil
		}
	}
	return fmt.Errorf("%v have not %v", user, album)
}

func (m *mongoAlbum) Stop() {

}

func GetAlbum() *mongoAlbum {
	return album
}
