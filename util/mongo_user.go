package util

import (
	"fmt"
	. "github.com/orange-jacky/albums/data"
	. "github.com/orange-jacky/albums/db"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
)

type mongoUser struct {
	Mongo *MongoClient
}

var (
	user      *mongoUser
	user_once sync.Once
)

func MongoUser() *mongoUser {
	user_once.Do(func() {
		user = &mongoUser{}
		if err := user.Init(); err != nil {
			log.Fatalln(err)
		}
	})
	return user
}

//对外方法,使用时,先init,再start,退出时stop
func (m *mongoUser) Init() error {
	conf := GetConfigure()

	mongo := &MongoClient{}
	mongo.Hosts = conf.Mongo.Hosts
	mongo.Database = conf.Mongo.User.Db
	mongo.Collection = conf.Mongo.User.Collection
	m.Mongo = mongo
	return nil
}

func (m *mongoUser) CheckUser(user, passwd string) (int, string) {
	if user == "" || passwd == "" {
		return -1, fmt.Sprintf("user or passwd empty")
	}

	if err := m.Mongo.Connect(); err != nil {
		return -2, fmt.Sprintf("connect %v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	result := &User{}
	q := bson.M{"user": user}
	c := m.Mongo.GetCollection()
	err := c.Find(q).One(result)
	if err != nil { //用户不存在
		u := User{User: user, Password: passwd}
		if err := m.Mongo.Insert(u); err != nil {
			return -3, fmt.Sprintf("create new user %v %v", user, err)
		} else {
			return NEW_USER, fmt.Sprintf("create new user %v success", user)
		}
	} else { //用户存在
		if result.Password == passwd {
			return 2, fmt.Sprintf("%v exist, and passwd is correct", user)
		} else {
			return -4, fmt.Sprintf("%v exist, but passwd is incorrect", user)
		}
	}
}

func (m *mongoUser) FindUser(user string) bool {
	if user == "" {
		return false
	}
	if err := m.Mongo.Connect(); err != nil {
		return false
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	q := bson.M{"user": user}

	c := m.Mongo.GetCollection()
	count, _ := c.Find(q).Count()
	if count > 0 {
		return true
	}
	return false
}

func (m *mongoUser) Stop() {

}

func GetUser() *mongoUser {
	return user
}
