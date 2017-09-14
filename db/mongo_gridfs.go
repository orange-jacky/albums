package db

import (
	"gopkg.in/mgo.v2"
)

type MongoGridfs struct {
	MongoClient
	GridFS *mgo.GridFS
}

func (m *MongoGridfs) C() {
	fs := "fs"
	if m.Collection != "" {
		fs = m.Collection
	}
	m.GridFS = m.Session.Db.GridFS(fs)
}

func (m *MongoGridfs) GetGridFs() *mgo.GridFS {
	return m.GridFS
}
