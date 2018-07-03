package transaction

import (
	"bytes"
	"chaoshen.com/goblock/utils"
	address "chaoshen.com/goblock/address"
)


type TXOutput struct {
	Value        int
	PubKeyHash []byte // PubKeyHash not address
}



func (txo *TXOutput) Lock(addr string){
	fullPayload:=utils.Base58Decode([]byte(addr))
	pubKeyHash:=fullPayload[0:len(fullPayload)-address.ADDRESS_CHECKSUM_LEN]
	txo.PubKeyHash=pubKeyHash
}

func (txo *TXOutput)IsLockedByKeyHash(pubKeyHash []byte) bool {
	return bytes.Equal(txo.PubKeyHash,pubKeyHash)
}

func NewTXOutput(value int, address string) *TXOutput{
	txo:=&TXOutput{Value:value}
	txo.Lock(address)
	return txo
}