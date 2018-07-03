package transaction

import (
	"bytes"
	"chaoshen.com/goblock/address"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey []byte
}
/*
func (txi *TXInput) CanUnlockWith(data []byte) bool {
	return bytes.Equal(txi.ScriptSig, data)
}
*/

func (txi *TXInput) CheckKeyHash(pubKeyHash []byte) bool{
	hash:=address.HashPublicKey(txi.PubKey)
	versionPubKeyHash := append([]byte{address.VERSION}, hash...)
	return bytes.Equal(versionPubKeyHash,pubKeyHash)
}