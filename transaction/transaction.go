package transaction

import (
	"bytes"
	"chaoshen.com/goblock/utils"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const SUBSIDY = 25

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig []byte
}

type TXOutput struct {
	Value        int
	ScriptPubKey []byte
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txInput := TXInput{[]byte{}, -1, []byte(data)}
	txOutput := TXOutput{SUBSIDY, []byte(to)}
	tx := Transaction{Vin: []TXInput{txInput}, Vout: []TXOutput{txOutput}}
	tx.SetID()
	return &tx
}

func (tx *Transaction) IsCoinbaseTx() bool {
	if len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1 {
		return true
	}
	return false
}

func (txi *TXInput) CanUnlockWith(data []byte) bool {
	return bytes.Equal(txi.ScriptSig, data)
}

func (txo *TXOutput) CanBeUnlockWith(data []byte) bool {
	return bytes.Equal(txo.ScriptPubKey, data)
}

func ValidateTransactions(trans []*Transaction) (remain []*Transaction, discard []*Transaction) {
	used := make(map[string][]int)
	for _, tran := range trans {
		flag:=false
	Loop:
		for _, in := range tran.Vin {
			if index, ok := used[string(in.Txid)]; ok {
				if utils.ContainInt(index, in.Vout) {
					discard = append(discard, tran)
					flag=true
					break Loop
				}
			}
		}
		if flag ==false {
			for _, in := range tran.Vin {
				used[string(in.Txid)] = append(used[string(in.Txid)], in.Vout)
			}
			remain = append(remain, tran)
		}
	}
	return remain,discard
}
