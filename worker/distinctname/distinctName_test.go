package main

import (
	"flag"
	"fmt"
	"strconv"
	"testing"

	"gopkg.in/redis.v3"
)

var eventTests = []struct {
	entry    string
	expected int
}{
	{"kite_call", 1},
	{"kite_something", 1},
	{"kite_call", 2},
	{"kite_call", 3},
}

func init() {
	flag.Parse()
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", *HOST, *PORT),
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	client.FlushDb()
}

func TestSetNewEvent(t *testing.T) {
	for _, tt := range eventTests {
		actual := SetEvent(tt.entry)
		if actual != float64(tt.expected) {
			t.Errorf("Set(%d): expected %d, actual %d", tt.entry, tt.expected, actual)
			fmt.Println("fail")
		}
	}
}

// Check that a new "bucket" is created and that the entries are cleared
func TestSetMonthlyBucket(t *testing.T) {
	client.ZIncr(strconv.FormatInt(DayOfUnixMonth()+1, 10), redis.Z{float64(111), "kite_call"})
	fmt.Println(client.DbSize().Result())
	SetMonthlyBucket()
	actual, _ := client.DbSize().Result()
	fmt.Println(actual)
	if actual != 1 {
		t.Error("The number of keys should be 0: actual is %d ", actual)
	}
}
