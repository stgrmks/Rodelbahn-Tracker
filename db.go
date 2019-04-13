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

func (d *RbData) Commit(c *mgo.Collection) error {
	if err := c.Insert(d); err != nil {
		log.Fatal("Insert failed: ", err)
		return err
	} else {
		log.Info("Insert success: ", d)
		return nil
	}
}

func Read(c mgo.Collection, f interface{}) ([]RbData, error) {
	result := make([]RbData, 0)
	if err := c.Find(f).All(&result); err != nil {
		log.Print("Read failed: ", err)
		return result, err
	} else {
		log.Print("Read success.")
		return result, err
	}
}

func InitDBConn(c *Config) {
	info := &mgo.DialInfo{
		Addrs:    []string{c.Host},
		Timeout:  60 * time.Second,
		Database: c.Database,
		Username: c.User,
		Password: c.Password,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	c.ActiveSession = session
	log.Info("Connection to DB established.")
}
