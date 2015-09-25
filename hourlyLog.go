package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"time"
)

type MetricData struct {
	Username string
	Count    int64
	Metric   string
	T        time.Time
}

var c *mgo.Collection

func main() {
	session := NewMongoClient()
	defer session.Close()
	Listen()
}

func NewMongoClient() *mgo.Session {
	session, err := mgo.Dial("192.168.99.100:32769")
	if err != nil {
		panic(err)
	}
	c = session.DB("metric").C("entries")
	return session
}

func MessageRead(body []byte) {
	res := MetricData{}
	json.Unmarshal(body, &res)
	fmt.Println(res)
	SetEvent(res)
}

func SetEvent(event MetricData) {
	err := c.Insert(&MetricData{event.Username, event.Count, event.Metric, time.Now()})
	if err != nil {
		log.Fatal(err)
	}
	index := mgo.Index{
		Key:         []string{"t"},
		ExpireAfter: 1 * time.Hour,
	}
	err = c.DropIndex("t")
	err = c.EnsureIndex(index)
	if err != nil {
		log.Fatal(err)
	}
}

func AverageOfValues() {
	//pipe := c.Pipe([]bson.M{{"$match": bson.M{"name": "Otavio"}}})
}
