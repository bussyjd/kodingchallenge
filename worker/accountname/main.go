package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	"database/sql"
	"encoding/json"
	"expvar"
	_ "github.com/lib/pq"
	"kodingchallenge/rabbit"
)

var (
	counts = expvar.NewMap("counters")
)

type MetricData struct {
	Username string
	Count    int64
	Metric   string
}

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "postgres"
)

var DEBUG = flag.Bool("debug_mode", false, "DEBUG MODE: ")
var HOST = flag.String("postgres_host", "127.0.0.1", "Postgres host (default 127.0.0.1): ")
var PORT = flag.Int("postgres_port", 5432, "Postgres port (default 5432): ")
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
	// Postgresql
	db := NewPsql()
	InitPsql(db)
	// Rabbitmq listener
	rabbit.Listen(*amqp_host, *amqp_port, counts, func(body []byte) {
		MessageRead(body, db)
	})
}

func NewPsql() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		*HOST, *PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if *DEBUG == true {
		fmt.Printf("postgresql connection info: %s \n", dbinfo)
	}
	checkErr(err)
	return db
}

func InitPsql(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS entry(username character varying(100) UNIQUE NOT NULL,Created date);")
	checkErr(err)
}

func WritePsql(username string, db *sql.DB) {
	statement, err := db.Prepare("INSERT INTO entry(username,created) VALUES($1,$2)")
	checkErr(err)
	_, err = statement.Exec(username, time.Now())
	checkErr(err)
	counts.Add("sql write", 1)
}

func MessageRead(body []byte, db *sql.DB) {
	res := MetricData{}
	json.Unmarshal(body, &res)
	WritePsql(res.Username, db)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
