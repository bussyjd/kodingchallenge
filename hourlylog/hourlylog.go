package hourlylog

import (
	"encoding/json"
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

func NewMongoClient(db string, collection string) *mgo.Session {
	if db == "" {
		db = "metric"
	}
	if collection == "" {
		collection = "entries"
	}
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
	err = c.DropIndex("t_1")
	err = c.EnsureIndex(index)
	if err != nil {
		log.Fatal(err)
	}
}

func CountEvent() int {
	n, err := c.Count()
	if err != nil {
		panic(err)
	}
	return n
}

func DropEventCollection() {
	err := c.DropCollection()
	if err != nil {
		panic(err)
	}
}
