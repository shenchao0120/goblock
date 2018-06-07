package core

import (
	"testing"
	"fmt"
)

func TestNewGenesisBlock(t *testing.T) {
	block:=NewGenesisBlock()
	fmt.Println(block.Header.Hash,string(block.Data))
}
