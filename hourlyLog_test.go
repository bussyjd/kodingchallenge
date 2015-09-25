package main

import (
	"./hourlyLog"
	//"fmt"
	//	"gopkg.in/mgo.v2"
	"testing"
	//"time"
)

func init() {
	hourlylog.NewMongoClient()
	//c = session.DB("metric").C("entries")
}

func TestMongoLogRecord(t *testing.T) {
	hourlylog.MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`))
	n := hourlylog.CountEvent()
	if n != 1 {
		t.Errorf("Expected 1 entry in the collection go %d", n)
	}
}

func TestMongoLogRecordExpire(t *testing.T) {
	//hourlylog.c.DropCollection()
	//err := c.Insert(&MetricData{"kodingbot_test", 12412415, "kite_call", time.Now().Add(-7200 * time.Second)})
	//if err != nil {
	//	panic(err)
	//}
	//MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`))

	//n, err := c.Count()
	//if n != 1 {
	//	t.Errorf("Expected 1 entry in the collection go %d", n)
	//}
}
