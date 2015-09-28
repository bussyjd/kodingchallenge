package accountname

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

type MetricData struct {
	Username string
	Count    int64
	Metric   string
}

const (
	HOST        = "192.168.99.100"
	PORT        = "32770"
	DB_USER     = "postgres"
	DB_PASSWORD = "postgres"
	DB_NAME     = "test"
)

func NewPsql() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
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
