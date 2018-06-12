package main

import (
	"chaoshen.com/goblock/core"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("New block.")

func main() {
	bc := core.NewBlockchain()
	defer bc.Close()
	err := bc.AddBlock([]byte("first message"))
	CheckErr(err)

	err=bc.AddBlock([]byte("second message"))
	CheckErr(err)
	err=bc.AddBlock([]byte("third message"))
	CheckErr(err)

	bc.Print()
}

func CheckErr(err error) {
	if err != nil {
		logger.Error(err)
	}
	return
}
