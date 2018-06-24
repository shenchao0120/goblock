package core

import (
	"fmt"
	"testing"
	"chaoshen.com/goblock/transaction"
)

func TestNewGenesisBlock(t *testing.T) {
	//block := NewGenesisBlock()
	//fmt.Println(block.Header.Hash, string(block.Data))
	b:=NewBlock([]byte("123456"),2,[]*transaction.Transaction{&transaction.Transaction{[]byte("abcde"),nil,nil}})

	fmt.Println("hash",string(b.HashTransactions()))
}
