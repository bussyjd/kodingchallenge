package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"encoding/json"
	"expvar"
	"gopkg.in/mgo.v2"
	"kodingchallenge/rabbit"
)

var (
	counts = expvar.NewMap("counters")
)

var DEBUG = flag.Bool("debug_mode", false, "DEBUG MODE: ")
var HOST = flag.String("mongo_host", "127.0.0.1", "mongo (default 127.0.0.1): ")
var PORT = flag.Int("mongo_port", 27017, "mongo port (default 27017): ")
var amqp_host = flag.String("amqp_host", "127.0.0.1", "RabbitMQ host (default 127.0.0.1): ")
var amqp_port = flag.Int("amqp_port", 5672, "RAbbitMQ port (default 5672): ")

func main() {
	// Flag parameters parsing
	flag.Parse()
	// Metrics server
	sock, err := net.Listen("tcp", "localhost:8123")
	checkErr(err)
	go func() {
		if *DEBUG == true {
			fmt.Println("Metrics server now available at localhost:8123/debug/vars")
		}
		http.Serve(sock, nil)
	}()
	// Mongodb
	session := NewMongoClient("", "")
	defer session.Close()
	// Rabbitmq listener
	rabbit.Listen(*amqp_host, *amqp_port, counts, func(body []byte) {
		MessageRead(body)
	})
}

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
	if *DEBUG == true {
		fmt.Printf("Connection to mongo server: %s:%v\n", *HOST, *PORT)
	}
	session, err := mgo.Dial(fmt.Sprintf("%s:%v", *HOST, *PORT))
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
