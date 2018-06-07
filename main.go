package main

import "chaoshen.com/goblock/core"

func main()  {
	bc:=core.NewBlockchain()
	bc.AddBlock([]byte("first message"))
	bc.AddBlock([]byte("second message"))
	bc.AddBlock([]byte("third message"))

	bc.Print()
}
