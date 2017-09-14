package db

import (
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

type Session struct {
	S  *mgo.Session
	Db *mgo.Database
	C  *mgo.Collection
}

//Mongo 定义mongo 操作实例
type MongoClient struct {
	Hosts      string
	Database   string
	Collection string
	Session
}

func (m *MongoClient) Connect() error {
	dailinfo := &mgo.DialInfo{
		Addrs:    strings.Split(m.Hosts, ","),
		Timeout:  5 * time.Second,
		Database: m.Database,
	}
	session, err := mgo.DialWithInfo(dailinfo)
	if err != nil {
		return err
	}
	m.S = session
	// Optional. Switch the session to a monotonic behavior.
	m.S.SetMode(mgo.Monotonic, true)
	return nil
}

func (m *MongoClient) DB() {
	db := m.Session.S.DB(m.Database)
	m.Session.Db = db
}

func (m *MongoClient) C() {
	c := m.Session.Db.C(m.Collection)
	m.Session.C = c
}

func (m *MongoClient) Close() {
	if m.S != nil {
		m.S.LogoutAll()
		m.S.Close()
	}
}

func (m *MongoClient) GetCollection() *mgo.Collection {
	return m.Session.C
}
