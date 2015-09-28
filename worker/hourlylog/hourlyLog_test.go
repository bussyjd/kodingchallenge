package main

import (
	//"fmt"
	_ "gopkg.in/mgo.v2"
	"testing"
	"time"
)

func init() {
	NewMongoClient("metric_test", "entries_test")
	DropEventCollection()
}

func TestMongoLogRecord(t *testing.T) {
	MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`))
	n := CountEvent()
	if n != 1 {
		t.Errorf("Expected 1 entry in the collection got %d", n)
	}
	defer DropEventCollection()
}

func TestMongoLogRecordExpire(t *testing.T) {
	SetEvent(MetricData{
		Username: "kodingbot_test",
		Count:    12412415,
		Metric:   "kite_call",
		T:        time.Now().Add(-7200 * time.Second)})
	MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`))

	n := CountEvent()
	if n != 1 {
		t.Errorf("Expected 1 entry in the collection go %d", n)
	}
	defer DropEventCollection()
}
