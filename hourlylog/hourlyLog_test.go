package main

import (
	"./hourlyLog"
	//"fmt"
	//	"gopkg.in/mgo.v2"
	"testing"
	"time"
)

var pastentry struct{}

func init() {
	hourlylog.NewMongoClient("metric_test", "entries_test")
}

func TestMongoLogRecord(t *testing.T) {
	hourlylog.MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`))
	n := hourlylog.CountEvent()
	if n != 1 {
		t.Errorf("Expected 1 entry in the collection got %d", n)
	}
	hourlylog.DropEventCollection()
}

func TestMongoLogRecordExpire(t *testing.T) {
	hourlylog.SetEvent(pastentry{
		Username: "kodingbot_test",
		Count:    12412415,
		Metric:   "kite_call",
		T:        time.Now().Add(-7200 * time.Second)})
	hourlylog.MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`))

	n := hourlylog.CountEvent()
	if n != 1 {
		t.Errorf("Expected 1 entry in the collection go %d", n)
	}
}
