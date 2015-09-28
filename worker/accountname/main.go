package main

import (
	"flag"
	"fmt"
	"time"

	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"kodingchallenge/rabbit"
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
	flag.Parse()
	db := NewPsql()
	InitPsql(db)
	rabbit.Listen(*amqp_host, *amqp_port, func(body []byte) {
		MessageRead(body, db)
	})
}

func NewPsql() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable",
		*HOST, *PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if *DEBUG == true {
		fmt.Printf("postgresql connction info: %s \n", dbinfo)
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
}

func MessageRead(body []byte, db *sql.DB) {
	res := MetricData{}
	json.Unmarshal(body, &res)
	fmt.Println(res)
	fmt.Println(res.Username)
	WritePsql(res.Username, db)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
