package core

import (
	"bytes"
	trans "chaoshen.com/goblock/transaction"
	"chaoshen.com/goblock/utils"
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"strings"
	"time"
	"log"
)

const dbFile = "blockchain.db"

const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

const Separator = "||"

const SELF_ADDRESS = "SELF_ADDRESS"

type Blockchain struct {
	Height int64
	dh     *DbHandler
}

// BlockchainIterator is used to iterate over blockchain blocks
type BlockchainIterator struct {
	currentHeight int64
	bc            *Blockchain
}

func NewBlockchain(address string) *Blockchain {
	dh := GetDbInstance(dbFile)
	height, err := dh.ResumeBlock()
	if err != nil {
		logger.Error(err)
		return nil
	}
	if height <= 0 {
		coinbase := trans.NewCoinbaseTx(address, "")
		err := dh.PutBlock(NewGenesisBlock(coinbase))
		height = 1
		if err != nil {
			logger.Error(err)
			return nil
		}
	}
	return &Blockchain{height, dh}
}

func CreateBlockchain(address string) *Blockchain {
	dh:=CreateDbInstance(dbFile)
	coinbase := trans.NewCoinbaseTx(address, "")
	err := dh.PutBlock(NewGenesisBlock(coinbase))
	height:= 1
	if err != nil {
		logger.Error(err)
		return nil
	}
	return &Blockchain{int64(height), dh}
}

func GetBlockchain() *Blockchain {
	dh := GetDbInstance(dbFile)
	height, err := dh.ResumeBlock()
	if err != nil {
		logger.Error(err)
		return nil
	}
	if  height<=0 {
		log.Panic(errors.New("The blockchain is not exist."))
		return nil
	}

	return &Blockchain{height, dh}
}
func (bc *Blockchain) GetBlock(height int64) *Block {
	return bc.dh.GetBlock(height)
}

func (bc *Blockchain) GenerateBlock(txs []*trans.Transaction) error {
	done := make(chan struct{})
	last := bc.GetBlock(bc.Height)
	coinbase := trans.NewCoinbaseTx(SELF_ADDRESS,"")
	txs=append(txs,coinbase)
	block := NewBlock(last.Header.Hash, last.Header.Height+1, txs)
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
	err := bc.dh.PutBlock(block)
	if err != nil {
		return err
	}
	bc.Height = block.Header.Height
	//bc.blocks = append(bc.blocks, block)
	return nil
}

// Iterator ...
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{1, bc}

	return bci
}

func (bci *BlockchainIterator) Next() *Block {
	block := bci.bc.GetBlock(bci.currentHeight)
	bci.currentHeight++
	return block
}
func (bci *BlockchainIterator) HasNext() bool {
	block := bci.bc.GetBlock(bci.currentHeight)
	if block != nil {
		return true
	}
	return false
}

func (bc *Blockchain) Close() {
	bc.dh.Db.Close()
}

func (bc *Blockchain) Print() {
	fmt.Println("Begin to print blockchain data:")
	bci := bc.Iterator()
	for bci.HasNext() {
		b := bci.Next()
		fmt.Printf("Heigth:%d,Prehash:%v,Hash:%v,data:%s\n", b.Header.Height, b.Header.PreBlockHash, b.Header.Hash, string(b.Transactions[0].ID))
	}
}

func (bc *Blockchain) FindUnspentTransactions(address string) map[*trans.Transaction][]int {
	unspentTXOs := make(map[string][]int)
	transMap := make(map[string]*trans.Transaction)

	unspentTrans := make(map[*trans.Transaction][]int)

	bci := bc.Iterator()
	for bci.HasNext() {
		block := bci.Next()
		for _, tx := range block.Transactions {
			transMap[string(tx.ID)] = tx
			for index, output := range tx.Vout {
				if output.CanBeUnlockWith([]byte(address)) {
					unspentTXOs[string(tx.ID)] = append(unspentTXOs[string(tx.ID)], index)
				}
			}
			if tx.IsCoinbaseTx() == false {
				for _, input := range tx.Vin {
					if idxs, ok := unspentTXOs[string(input.Txid)]; ok && input.CanUnlockWith([]byte(address)) {
						unspentTXOs[string(input.Txid)] = utils.RemoveIntInSlice(idxs, input.Vout)
					}
				}
			}
		}
	}
	for key, val := range unspentTXOs {
		if len(val) == 0 {
			delete(unspentTXOs, key)
			continue
		}
		tran := transMap[key]
		unspentTrans[tran] = val
	}
	return unspentTrans
}

func (bc *Blockchain) FindUTXOs(address string) []*trans.TXOutput {
	var txOutput []*trans.TXOutput
	transMap := bc.FindUnspentTransactions(address)
	for trans, indexs := range transMap {
		for _, i := range indexs {
			txOutput = append(txOutput, &trans.Vout[i])
		}
	}
	return txOutput
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string]*trans.TXOutput) {
	res := make(map[string]*trans.TXOutput)
	transMap := bc.FindUnspentTransactions(address)
	sum := 0
Find:
	for tran, indexs := range transMap {
		for _, i := range indexs {
			txOut := tran.Vout[i]
			key := bytes.Join([][]byte{tran.ID, []byte(strconv.Itoa(i))}, []byte(Separator))
			if txOut.Value >= amount {
				res = make(map[string]*trans.TXOutput)
				res[string(key)] = &txOut
				sum = txOut.Value
				break Find
			}
			sum += txOut.Value
			res[string(key)] = &txOut
			if sum >= amount {
				break Find
			}
		}
	}
	return sum, res
}

func (bc *Blockchain) NewUTXOTransaction(from, to string, amount int) (*trans.Transaction, error) {
	sumAmt, outputs := bc.FindSpendableOutputs(from, amount)
	if sumAmt < amount {
		return nil, errors.New("Not enough funds.")
	}
	remain := sumAmt - amount
	tx := trans.Transaction{}
	// build input
	for key, _ := range outputs {
		keys := strings.Split(key, Separator)
		vout, err := strconv.Atoi(keys[1])
		if err != nil {
			logs.Error(err)
		}
		input := trans.TXInput{[]byte(keys[0]), vout, []byte(from)}
		tx.Vin = append(tx.Vin, input)
	}
	//build output
	toOutput := trans.TXOutput{amount, []byte(to)}
	tx.Vout = append(tx.Vout, toOutput)
	if remain >0 {
		remainOutput := trans.TXOutput{remain, []byte(from)}
		tx.Vout = append(tx.Vout, remainOutput)
	}
	tx.SetID()
	return &tx,nil
}

func (bc *Blockchain) GetBalance(address string) int {
	utxos:=bc.FindUTXOs(address)
	balance:=0

	for _,utxo := range utxos {
		balance+=utxo.Value
	}
	return balance
}


func (bc *Blockchain) GetTransactionNum()int {
	bci:=bc.Iterator()
	tranNum:=0
	for bci.HasNext(){
		block:=bci.Next()
		tranNum+=len(block.Transactions)
	}
	return tranNum
}
