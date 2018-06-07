package core

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type Block struct {
	Header BlockHeader
	Data   []byte
}

type BlockHeader struct {
	PreBlockHash []byte
	Hash         []byte
	Nonce        uint32
	Height       int64
	Timestamp    int64
}

func NewBlock(preHash []byte, nonce uint32, height int64, data []byte) *Block {
	timestamp:=time.Now().Unix()
	block:=&Block{BlockHeader{PreBlockHash:preHash,Nonce:nonce,Height:height,Timestamp:timestamp},data}
	block.SetHash()
	return block
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Header.Timestamp, 10))
	nonce := []byte(strconv.FormatUint(uint64(b.Header.Nonce), 10))
	height := []byte(strconv.FormatInt(b.Header.Height, 10))

	head := bytes.Join([][]byte{b.Header.PreBlockHash, nonce, height, b.Data, timestamp}, []byte{})
	//h := sha256.New()
	//h.Write(head)
	//b.Header.Hash = h.Sum(nil)
	h:=sha256.Sum256(head)
	b.Header.Hash=h[:]
}

func NewGenesisBlock() *Block {
	return NewBlock([]byte{},0,1,[]byte("Genesis block"))
}
