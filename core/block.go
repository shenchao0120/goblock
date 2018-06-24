package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/op/go-logging"
	"strconv"
	"time"
	. "chaoshen.com/goblock/transaction"
)

var logger = logging.MustGetLogger("Block")

type Block struct {
	Header BlockHeader
	Transactions [] *Transaction
}

type BlockHeader struct {
	PreBlockHash []byte
	Hash         []byte
	Nonce        int64
	Height       int64
	Timestamp    int64
}

func NewBlock(preHash []byte, height int64, trans []*Transaction) *Block {
	timestamp := time.Now().Unix()
	block := &Block{BlockHeader{PreBlockHash: preHash, Height: height, Timestamp: timestamp}, trans}
	block.SetHash()
	return block
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Header.Timestamp, 10))
	nonce := []byte(strconv.FormatUint(uint64(b.Header.Nonce), 10))
	height := []byte(strconv.FormatInt(b.Header.Height, 10))

	head := bytes.Join([][]byte{b.Header.PreBlockHash, nonce, height, b.HashTransactions(), timestamp}, []byte{})
	//h := sha256.New()
	//h.Write(head)
	//b.Header.Hash = h.Sum(nil)
	h := sha256.Sum256(head)
	b.Header.Hash = h[:]
}

func (b *Block)HashTransactions()[]byte{
	var hash []byte
	for _,tx:=range b.Transactions{
		hash=append(hash,tx.ID...)
	}
	result:=sha256.Sum256(hash)
	return result[:]
}

func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(b)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func NewGenesisBlock(CoinbaseTx *Transaction) *Block {
	return NewBlock([]byte{}, 1, []*Transaction{CoinbaseTx})
}

func DeserializeBlock(b []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(b))
	err := decoder.Decode(&block)
	if err != nil {
		logger.Error("Error:", err)
	}
	return &block
}
