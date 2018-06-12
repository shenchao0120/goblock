package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

const dbFile = "blockchain.db"

type Blockchain struct {
	Height int64
	dh     *DbHandler
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHeight int64
	bc          *Blockchain
}

func NewBlockchain() *Blockchain {
	dh := GetDbInstance(dbFile)
	height, err := dh.ResumeBlock()
	if err != nil {
		logger.Error(err)
		return nil
	}
	if height <= 0 {
		err := dh.PutBlock(NewGenesisBlock())
		height = 1
		if err != nil {
			logger.Error(err)
			return nil
		}
	}
	return &Blockchain{height,dh}
}

func (bc *Blockchain) GetBlock(height int64) *Block {
	return bc.dh.GetBlock(height)
}

func (bc *Blockchain) AddBlock(data []byte) error {
	done := make(chan struct{})
	last := bc.GetBlock(bc.Height)
	block := NewBlock(last.Header.Hash, last.Header.Height+1, data)
	pow := NewProofOfWork(block, 18)
	ctx, cancle := context.WithTimeout(context.Background(), 60*time.Second)
	ticker := time.After(60 * time.Second)
	defer cancle()
	go pow.Run(ctx, done)

	select {
	case <-done:
		fmt.Println("Receive done channel.")
	case <-ticker:
		fmt.Println("Receive Timeout event.")
	}
	if block.Header.Nonce == 0 {
		logs.Error("mint nonce error")
		return errors.New("Mint nonce error.")
	}
	err:=bc.dh.PutBlock(block)
	if err != nil {
		return err
	}
	bc.Height=block.Header.Height
	//bc.blocks = append(bc.blocks, block)
	return nil
}

// Iterator ...
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{1, bc}

	return bci
}


func (bci *BlockchainIterator)Next() *Block{
	block:=bci.bc.GetBlock(bci.currentHeight)
	bci.currentHeight++
	return block
}
func (bci *BlockchainIterator)HasNext() bool{
	block:=bci.bc.GetBlock(bci.currentHeight)
	if block != nil{
		return true
	}
	return false
}

func (bc *Blockchain) Close()  {
	bc.dh.Db.Close()
}


func (bc *Blockchain) Print() {
	bci:=bc.Iterator()
	for bci.HasNext() {
		b:=bci.Next()
		fmt.Printf("Heigth:%d,Prehash:%v,Hash:%v,data:%s\n", b.Header.Height, b.Header.PreBlockHash, b.Header.Hash, string(b.Data))
	}
}


