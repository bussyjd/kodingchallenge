package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"expvar"
	"gopkg.in/redis.v3"
	"kodingchallenge/rabbit"
)

var (
	counts = expvar.NewMap("counters")
)

const key = "metric"

var client *redis.Client

type MetricData struct {
	Username string
	Count    int64
	Metric   string
}

var DEBUG = flag.Bool("debug_mode", false, "DEBUG MODE: ")
var HOST = flag.String("redis_host", "127.0.0.1", "Redis (default 127.0.0.1): ")
var PORT = flag.Int("redis_port", 6379, "Redis port (default 6379): ")
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
	// Rabbitmq listener
	go rabbit.Listen(*amqp_host, *amqp_port, counts, func(body []byte) {
		MessageRead(body)
	})
	// Redis
	NewClient()
	BucketCheck()
}

func NewClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", *HOST, *PORT),
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	if *DEBUG == true {
		fmt.Println(pong, err)
	}
}

func BucketCheck() {
	for range time.Tick(time.Minute) {
		if *DEBUG == true {
			fmt.Println("tick")
		}
		if DayOfUnixMonth == nil {
			exists, err := client.Exists(strconv.FormatInt(UnixTimeHour(), 10)).Result()
			if exists != true {
				SetMonthlyBucket()
			}
			if err != nil {
				fmt.Errorf("Eror while Checking for existing bucket: %d", err)
			}
		}
	}
}

func MessageRead(body []byte) {
	res := MetricData{}
	json.Unmarshal(body, &res)
	SetEvent(res.Metric)
}

func SetEvent(name string) float64 {
	fmt.Println(DayOfUnixMonth())
	result, err := client.ZIncr(strconv.FormatInt(DayOfUnixMonth(), 10), redis.Z{float64(1), name}).Result()

	if *DEBUG == true {
		fmt.Println(result, err)
	}
	return result
}

func SetMonthlyBucket() int64 {
	a := make([]string, 30, 30)
	for i, _ := range a {
		a[i] = strconv.Itoa(i + 1)
	}
	if *DEBUG == true {
		fmt.Printf("%v\n", a)
	}
	merged, err := client.ZUnionStore(
		strconv.FormatInt(time.Now().Unix(), 10),
		redis.ZStore{Aggregate: "SUM"},
		a...,
	).Result()
	if err == nil {
		client.Del(a...)
	}
	return merged
}

// the count occurs every 30 days, not monthly so unix time is used as a reference
func DayOfUnixMonth() int64 {
	unixtimeday := UnixTimeHour()
	return unixtimeday % 30
}

func UnixTimeHour() int64 {
	unixtime := time.Now().Unix()
	return (((unixtime / 60) / 60) / 24)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
