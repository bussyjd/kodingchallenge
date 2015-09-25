package main

import (
	"./accountName"
	"./rabbit"
)

func main() {
	db := accountname.NewPsql()
	accountname.InitPsql(db)
	rabbit.Listen(func(body []byte) {
		accountname.MessageRead(body, db)
	})
}
