package distinctname

import (
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"strconv"
	"time"
)

const key = "metric"

var client *redis.Client

type MetricData struct {
	Username string
	Count    int64
	Metric   string
}

func BucketCheck() {
	for range time.Tick(time.Minute) {
		fmt.Println("tick")
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
	//fmt.Printf(" [x] %s", body)
	res := MetricData{}
	json.Unmarshal(body, &res)
	fmt.Println(res)
	SetEvent(res.Metric)
}

func NewClient() {
	client = redis.NewClient(&redis.Options{
		Addr:     "192.168.99.100:32768",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func SetEvent(name string) float64 {
	fmt.Println(DayOfUnixMonth())
	result, err := client.ZIncr(strconv.FormatInt(DayOfUnixMonth(), 10), redis.Z{float64(1), name}).Result()
	fmt.Println(result, err)
	return result
}

func SetMonthlyBucket() int64 {
	a := make([]string, 30, 30)
	for i, _ := range a {
		a[i] = strconv.Itoa(i + 1)
	}
	fmt.Printf("%v\n", a)
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
