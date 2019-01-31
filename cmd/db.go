package cmd

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
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

func RunInitDB() {
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
	config.ActiveSession = session

	fmt.Println("MongoDB connection established.")

	config.ActiveCollection = session.DB(config.Database).C(config.Collection)

	index := mgo.Index{
		Key:    []string{"time", "location", "user", "rating", "comment"},
		Unique: true,
	}

	if err := config.ActiveCollection.EnsureIndex(index); err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v\n", config)
}

func (d *RbData) Commit(c *mgo.Collection) error {
	if err := c.Insert(d); err != nil {
		log.Print("Insert failed: ", err)
		return err
	} else {
		log.Print("Insert success: ", d)
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
