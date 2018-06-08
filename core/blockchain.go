package core

import (
	"fmt"
	"context"
	"time"
	"github.com/astaxie/beego/logs"
	"qiniupkg.com/x/errors.v7"
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

func (bc *Blockchain) AddBlock(data []byte) error {
	done:=make(chan struct{})
	last:=bc.blocks[len(bc.blocks)-1]
	block:=NewBlock(last.Header.Hash,last.Header.Height+1,data)
	pow:=NewProofOfWork(block,18)
	ctx,cancle:=context.WithTimeout(context.Background(),60 * time.Second)
	defer cancle();
	go pow.Run(ctx,done)
	<-done
	if block.Header.Nonce == 0 {
		logs.Error("mint nonce error")
		return errors.New("Mint nonce error.")
	}
	bc.blocks =append(bc.blocks,block)
	return nil
}


func (bc *Blockchain) Print(){
	for _,b:=range bc.blocks {
		fmt.Printf("Heigth:%d,Prehash:%v,Hash:%v,data:%s\n",b.Header.Height,b.Header.PreBlockHash,b.Header.Hash,string(b.Data))
	}
}