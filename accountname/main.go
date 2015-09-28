package main

import (
	"kodingchallenge/accountname"
	"kodingchallenge/rabbit"
)

func main() {
	db := accountname.NewPsql()
	accountname.InitPsql(db)
	rabbit.Listen(func(body []byte) {
		accountname.MessageRead(body, db)
	})
}
