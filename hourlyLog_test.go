package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"testing"
	"time"
)

//var c *mgo.Collection

func init() {
	session, err := mgo.Dial("192.168.99.100:32769")
	if err != nil {
		panic(err)
	}
	c = session.DB("metric_test").C("entries")
	c.DropCollection()
}

func TestMongoLogRecord(t *testing.T) {
	MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`))
	n, err := c.Count()
	if err != nil {
		panic(err)
	}
	if n != 1 {
		t.Errorf("Expected 1 entry in the collection go %d", n)
	}
}

func TestMongoLogRecordExpire(t *testing.T) {
	c.DropCollection()
	err := c.Insert(&MetricData{"kodingbot_test", 12412415, "kite_call", time.Now().Add(-7200 * time.Second)})
	if err != nil {
		panic(err)
	}
	MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`))

	n, err := c.Count()
	if n != 1 {
		t.Errorf("Expected 1 entry in the collection go %d", n)
	}
}
