package main

import (
	"kodingchallenge/hourlylog"
	"kodingchallenge/rabbit"
)

func main() {
	session := hourlylog.NewMongoClient()
	defer session.Close()
	rabbit.Listen(func(body []byte) {
		hourlylog.MessageRead(body)
	})
}
