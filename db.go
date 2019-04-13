package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type RbData struct {
	Id       bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Time     time.Time     `bson:"time" bson:"time"`
	Location string        `bson:"location" bson:"location"`
	Link     string        `bson:"Link" bson:"Link"`
	User     string        `bson:"User" bson:"User"`
	Rating   string        `bson:"rating" bson:"rating"`
	Comment  string        `bson:"comment" bson:"comment"`
}

type DbSession struct {
	Session    *mgo.Session
	Collection *mgo.Collection
	Schema     RbData
}

func (dbs *DbSession) Connect(config *Config) {
	info := &mgo.DialInfo{
		Addrs:    []string{config.Host},
		Timeout:  60 * time.Second,
		Database: config.Database,
		Username: config.User,
		Password: config.Password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	dbs.Session = session
	log.Infof("Connection to DB: %s established.", config.Database)

	dbs.Collection = session.DB(config.Database).C(config.Collection)
	log.Info("Collection: %s selected.", config.Collection)

	// allow only unique entries by using unique index generated by data
	index := mgo.Index{
		Key:    []string{"time", "location", "user", "rating", "comment"},
		Unique: true,
	}

	if err := dbs.Collection.EnsureIndex(index); err != nil {
		log.Fatal(err)
	}
	log.Infof("%+v\n", dbs)
}

func (dbs *DbSession) Disconnect() {
	dbs.Session.Close()
	log.Info("Disconnected from DB.")
}

func (dbs *DbSession) Commit(d []RbData) {

	defer dbs.Session.Close()

	for _, tmpData := range d {
		if err := dbs.Collection.Insert(tmpData); err != nil {
			log.Debugf("Insert failed: %s", err)
		} else {
			log.Debug("Insert success: %s", tmpData)
		}
	}
}
