package core
/*
import (
	"fmt"
	"strconv"
	"testing"
)

func TestDbHandler_PutBlock(t *testing.T) {
	dh := GetDbInstance(dbFile)
	if dh == nil {
		t.Error("nil db handler.")
	}
	isExist := dh.BlockChainIsExist()
	fmt.Println("blockchain exist:", isExist)
	gensis := NewGenesisBlock()
	fmt.Println(gensis)
	err := dh.PutBlock(gensis)
	if err != nil {
		fmt.Println(err)
	}
	block := dh.GetBlock(gensis.Header.Height)
	isExist = dh.BlockChainIsExist()
	fmt.Println("blockchain exist:", isExist)

	block = NewBlock(block.Header.Hash, block.Header.Height+1, nil)
	dh.PutBlock(block)
	fmt.Println(block)

	hash, _ := dh.ResumeBlock()
	heigth, _ := strconv.Atoi(string(hash))
	fmt.Println("the heigth:", heigth)
}
*/