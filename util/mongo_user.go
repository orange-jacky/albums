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

func (m *mongoUser) NewUser(user, passwd string) error {
	if err := m.Mongo.Connect(); err != nil {
		return fmt.Errorf("connect db fail,%v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	new_user := User{User: user, Password: passwd}
	if err := collection.Insert(new_user); err != nil {
		return fmt.Errorf("create new user %v %v", user, err)
	}
	return nil
}

func (m *mongoUser) CheckUserAndPasswd(user, passwd string) (int, error) {
	if err := m.Mongo.Connect(); err != nil {
		return -1, fmt.Errorf("connect db fail,%v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()

	collection := m.Mongo.GetCollection()

	q := bson.M{"user": user, "password": passwd}
	count, err := collection.Find(q).Count()

	if err != nil {
		return -1, fmt.Errorf("find %v fail,%v", user)
	}
	return count, nil
}

func (m *mongoUser) CheckUser(user string) (int, error) {
	if err := m.Mongo.Connect(); err != nil {
		return -1, fmt.Errorf("connect db fail,%v", err)
	}
	defer m.Mongo.Close()
	m.Mongo.DB()
	m.Mongo.C()
	collection := m.Mongo.GetCollection()

	q := bson.M{"user": user}
	count, err := collection.Find(q).Count()
	if err != nil {
		return -1, fmt.Errorf("find %v fail,%v", user)
	}
	return count, nil
}

func (m *mongoUser) Stop() {

}

func GetUser() *mongoUser {
	return user
}
