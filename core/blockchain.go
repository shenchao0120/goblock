package core

import (
	"fmt"
	"math/rand"
)

type Blockchain struct {
	blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

/*
func (bc *Blockchain) AddBlock(b *Block) {
	last := bc.blocks[len(bc.blocks)-1]
	b.Header.Hash = last.Header.Hash
	b.SetHash()
	bc.blocks = append(bc.blocks, b)
}
*/

func (bc *Blockchain) AddBlockWithNonce(nonce uint32,data []byte) {
	last:=bc.blocks[len(bc.blocks)-1]
	block:=NewBlock(last.Header.Hash,nonce,last.Header.Height+1,data)
	bc.blocks =append(bc.blocks,block)
}

func (bc *Blockchain) AddBlock(data []byte) {
	nonce:=rand.Uint32()
	bc.AddBlockWithNonce(nonce,data)
}


func (bc *Blockchain) Print(){
	for _,b:=range bc.blocks {
		fmt.Printf("Heigth:%d,Prehash:%v,Hash:%v,data:%s\n",b.Header.Height,b.Header.PreBlockHash,b.Header.Hash,string(b.Data))
	}
}