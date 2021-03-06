package main

import (
	"flag"
	"fmt"
	"testing"

	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	flag.Parse()
	db = NewPsql()
	InitPsql(db)
}

func TestCollect(t *testing.T) {
	MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`), db)
	c := CountRows("kodingbot_test")
	if c != 1 {
		fmt.Println("Failed to write to the DB actual %s", c)
		panic("n db write")
	}
}

func TestCollectUniq(t *testing.T) {
	MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412414, "metric": "kite_call"}`), db)
	MessageRead([]byte(`{"username": "kodingbot_test", "count": 12412415, "metric": "kite_call"}`), db)
	c := CountRows("kodingbot_test")
	if c != 1 {
		fmt.Println("Failed to have only one entry per user %s", c)
		panic("n db write")
	}
}

func CountRows(username string) int {
	var count int
	statement, err := db.Prepare("SELECT COUNT (\"username\") FROM entry WHERE username=$1")
	checkErr(err)
	row := statement.QueryRow(username)
	err = row.Scan(&count)
	return count
}
