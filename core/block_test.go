package core

import (
	"fmt"
	"testing"
)

func TestNewGenesisBlock(t *testing.T) {
	block := NewGenesisBlock()
	fmt.Println(block.Header.Hash, string(block.Data))
}
