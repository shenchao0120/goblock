package transaction

import (
	"bytes"
	"chaoshen.com/goblock/utils"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const SUBSIDY = 25

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) SetID() {
	tx.ID = nil
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

	txInput := TXInput{[]byte{}, -1, nil, []byte(data)}
	//txOutput := TXOutput{SUBSIDY, []byte(to)}
	txOutput := NewTXOutput(SUBSIDY,to)
	tx := Transaction{Vin: []TXInput{txInput}, Vout: []TXOutput{*txOutput}}
	tx.SetID()
	return &tx
}

func (tx *Transaction) IsCoinbaseTx() bool {
	if len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1 {
		return true
	}
	return false
}

func ValidateTransactions(trans []*Transaction) (remain []*Transaction, discard []*Transaction) {
	used := make(map[string][]int)
	for _, tran := range trans {
		flag := false
	Loop:
		for _, in := range tran.Vin {
			if index, ok := used[string(in.Txid)]; ok {
				if utils.ContainInt(index, in.Vout) {
					discard = append(discard, tran)
					flag = true
					break Loop
				}
			}
		}
		if flag == false {
			for _, in := range tran.Vin {
				used[string(in.Txid)] = append(used[string(in.Txid)], in.Vout)
			}
			remain = append(remain, tran)
		}
	}
	return remain, discard
}

func (tx *Transaction) TrimmedCopy() *Transaction {
	var input []TXInput
	var output []TXOutput

	for _, in := range tx.Vin {
		input = append(input, TXInput{in.Txid, in.Vout, nil, nil})
	}

	copy(output, tx.Vout)
	return &Transaction{tx.ID, input, output}
}

func (tx *Transaction) String() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	lines = append(lines, fmt.Sprintf("----- Transaction input:"))
	for id, in := range tx.Vin {
		lines = append(lines, fmt.Sprintf("       Input :%d", id+1))
		lines = append(lines, fmt.Sprintf("         TXID :%x", in.Txid))
		lines = append(lines, fmt.Sprintf("         Vout :%d", in.Vout))
		lines = append(lines, fmt.Sprintf("         Sig :%x", in.Signature))
		lines = append(lines, fmt.Sprintf("         PubKey :%s", in.PubKey))
	}
	lines = append(lines, fmt.Sprintf("----- Transaction output:"))

	for id, out := range tx.Vout {
		lines = append(lines, fmt.Sprintf("       Output :%d", id+1))
		lines = append(lines, fmt.Sprintf("         Value :%d", out.Value))
		lines = append(lines, fmt.Sprintf("         PubKeyHash :%s\n", out.PubKeyHash))
	}
	return strings.Join(lines, "\n")
}

func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey, prevTXs map[string]*TXOutput) {
	if tx.IsCoinbaseTx() {
		return
	}
	txCopy := tx.TrimmedCopy()

	for inId, in := range txCopy.Vin {
		key := utils.MakeCompositeKey(in.Txid, in.Vout)
		out, ok := prevTXs[key]
		if !ok {
			log.Panic("Cannot find right tx output:", in.Txid, ",", in.Vout)
		}
		// Notice: use PubKeyHash as pubKey here.
		txCopy.Vin[inId].PubKey = out.PubKeyHash
		txCopy.SetID()
		txCopy.Vin[inId].PubKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, privKey, txCopy.ID)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inId].Signature = signature
	}
}

func (tx *Transaction) Verify(prevTXs map[string]*TXOutput) bool {
	if tx.IsCoinbaseTx() {
		return true
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inId, in := range txCopy.Vin {
		key := utils.MakeCompositeKey(in.Txid, inId)
		out, ok := prevTXs[key]
		if !ok {
			log.Panic("Cannot find right tx output:", in.Txid, ",", in.Vout)
		}
		txCopy.Vin[inId].PubKey = out.PubKeyHash
		txCopy.SetID()
		txCopy.Vin[inId].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		r.SetBytes(txCopy.ID[:len(txCopy.ID)/2])
		s.SetBytes(txCopy.ID[len(txCopy.ID)/2:])

		x := big.Int{}
		y := big.Int{}
		x.SetBytes(tx.Vin[inId].PubKey[:len(tx.Vin[inId].PubKey)/2])
		y.SetBytes(tx.Vin[inId].PubKey[len(tx.Vin[inId].PubKey)/2:])
		rawPubKey := &ecdsa.PublicKey{curve, &x, &y}
		if !ecdsa.Verify(rawPubKey, tx.Vin[inId].Signature, &r, &s) {
			return false
		}
	}
	return true
}
