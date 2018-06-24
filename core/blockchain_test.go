package core

import (
	"testing"
	"fmt"
	"chaoshen.com/goblock/transaction"
	"os"
	"github.com/astaxie/beego/logs"
)



func TestBlockchain_FindUnspentTransactions(t *testing.T) {

	_,err:=os.Stat("blockchain.db")
	if err != os.ErrNotExist {
		os.Remove("blockchain.db")
	}

	bc:=NewBlockchain(SELF_ADDRESS)
	b:=bc.GetBlock(1)

	txInput:=transaction.TXInput{b.Transactions[0].ID,0,[]byte("SELF_ADDRESS")}
	txOutput1:=transaction.TXOutput{10,[]byte("b")}
	txOutput2:=transaction.TXOutput{10,[]byte("c")}

	tx1:=&transaction.Transaction{[]byte("22222"),[]transaction.TXInput{txInput},[]transaction.TXOutput{txOutput1,txOutput2}}
	err=bc.GenerateBlock([]*transaction.Transaction{tx1})
	if err != nil {
		fmt.Println(err)
	}

	txInput1:=transaction.TXInput{[]byte("22222"),1,[]byte("c")}
	txOutput3:=transaction.TXOutput{5,[]byte("b")}
	txOutput4:=transaction.TXOutput{5,[]byte("b")}
	tx2:=&transaction.Transaction{[]byte("33333"),[]transaction.TXInput{txInput1},[]transaction.TXOutput{txOutput3,txOutput4}}

	err=bc.GenerateBlock([]*transaction.Transaction{tx2})
	if err != nil {
		fmt.Println(err)
	}

	//m:=bc.FindUnspentTransactions("b")
	//for k,v:= range m {
	//	fmt.Println(k)
	//	fmt.Println(v)
	//}

	sum,res:=bc.FindSpendableOutputs("b",50)


	fmt.Println(sum)

	for k,v:= range res {
		fmt.Println(k)
		fmt.Println(v)
	}

}

func TestBlockchain_NewUTXOTransaction(t *testing.T) {
	_,err:=os.Stat("blockchain.db")
	if err != os.ErrNotExist {
		os.Remove("blockchain.db")
	}

	bc:=NewBlockchain(SELF_ADDRESS)
	tx,err:=bc.NewUTXOTransaction(SELF_ADDRESS,"aaaa",12)
	if err != nil {
		logs.Error(err)
	}

	tx2,err:=bc.NewUTXOTransaction(SELF_ADDRESS,"aaaa",10)
	if err != nil {
		logs.Error(err)
	}

	fmt.Println(tx)
	remain,discard:=transaction.ValidateTransactions([]*transaction.Transaction{tx,tx2})
	fmt.Println(len(remain),len(discard))
	err=bc.GenerateBlock(remain)
	if err != nil {
		logs.Error(err)
	}
	//bc.Print()
	sum,res:=bc.FindSpendableOutputs("aaaa",10)


	fmt.Println("sum:",sum,"res:",len(res))

	for k,v:= range res {
		fmt.Println(k)
		fmt.Println(v)
	}

}

func TestDelete(t *testing.T){
	_,err:=os.Stat("blockchain.db")
	if err != os.ErrNotExist {
		os.Remove("blockchain.db")
	}
}