package db

import (
	. "github.com/orange-jacky/albums/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

//Mongo 定义mongo 操作实例
type Mongo struct {
	Hosts      string
	Database   string
	Collection string
	S          *mgo.Session
	Db         *mgo.Database
	C          *mgo.Collection
}

func NewMongo() *Mongo {
	return &Mongo{}
}

func (m *Mongo) Connect(hosts, db string) error {
	dailinfo := &mgo.DialInfo{
		Addrs:    strings.Split(hosts, ","),
		Timeout:  5 * time.Second,
		Database: db,
	}
	session, err := mgo.DialWithInfo(dailinfo)
	if err != nil {
		return err
	}
	m.S = session
	m.Database = db
	m.Hosts = hosts
	// Optional. Switch the session to a monotonic behavior.
	m.S.SetMode(mgo.Monotonic, true)
	return nil
}

func (m *Mongo) OpenDb(db string) error {
	m.Db = m.S.DB(db)
	return nil
}

func (m *Mongo) OpenTable(table string) error {
	m.C = m.Db.C(table)
	return nil
}

func (m *Mongo) Query(query interface{}) (features Features, err error) {
	record := new(Featuredata)
	iter := m.C.Find(query).Iter()
	for iter.Next(record) {
		features = append(features, record)
		record = new(Featuredata) //重置变量
	}
	if err := iter.Close(); err != nil {
		return features, err
	}
	return features, nil
}

func (m *Mongo) Insert(docs ...interface{}) error {
	return m.C.Insert(docs...)
}

func (m *Mongo) FindUserOne(username string) bool {
	count, err := m.C.Find(bson.M{"username": username}).Count()
	if err != nil {
		return false
	}
	if count != 0 {
		return true
	} else {
		return false
	}
}

func (m *Mongo) Close() error {
	if m.S != nil {
		m.S.LogoutAll()
		m.S.Close()
	}
	return nil
}
