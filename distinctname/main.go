package main

import (
	"kodingchallenge/distinctname"
	"kodingchallenge/rabbit"
)

func main() {
	go distincname.Listen()
	distincname.NewClient()
	distincname.BucketCheck()
}
