package core

import (
	"bytes"
	"chaoshen.com/goblock/utils"
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var MaxNonce = math.MaxInt64

const MaxBits = 256

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block, targetInt uint) *ProofOfWork {
	target := big.NewInt(1)

	target.Lsh(target, MaxBits-targetInt)
	return &ProofOfWork{block, target}
}

func (pow *ProofOfWork) PrepareData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			utils.ConvIntToHex(pow.block.Header.Height),
			utils.ConvIntToHex(pow.block.Header.Timestamp),
			pow.block.Header.PreBlockHash,
			pow.block.Data,
			utils.ConvIntToHex(nonce),
		}, []byte(""))
	return data
}

func (pow *ProofOfWork) Run(ctx context.Context, done chan<- struct{}) {
	nonce := 1
	var hash [32]byte
	var hashInt big.Int
	for nonce < MaxNonce {
		select {
		case <-ctx.Done():
			break
		default:
		}
		data := pow.PrepareData(int64(nonce))
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.target) == -1 {
			pow.block.Header.Nonce = int64(nonce)
			pow.block.Header.Hash = hash[:]
			fmt.Printf("Mine success nonce:%d,hash:%x\n", nonce, hash)
			done <- struct{}{}
			return
		}
		//if res:=bytes.Compare(hash[:],pow.target.Bytes()) ;res == -1 {
		//	return int64(nonce),hash
		//}
		//fmt.Printf("Mine failed nonce:%d,hash:%d,target:%d\n", nonce, &hashInt, pow.target)
		nonce++
		continue
	}
	done <- struct{}{}
	return
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.PrepareData(pow.block.Header.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])
	if hashInt.Cmp(pow.target) == -1 {
		return true
	}
	return false

}
