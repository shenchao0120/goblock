package address

import (
	"testing"
	"fmt"
)

func TestNewWallet(t *testing.T) {
	wallet:=NewWallet()
	address:=wallet.GetAddress()
	fmt.Printf("The address[%s],len[%d]\n",address,len(address))
	isValid:=ValidateAddress(address)

	fmt.Printf("Isvaild %v\n",isValid)
}



func TestNewWallets(t *testing.T) {
	ws,err:=NewWallets(wallet_file)
	if err != nil {
		t.Error(err)
	}
	add:=ws.CreateWallet()
	fmt.Printf("create address:%s",add)
	add=ws.CreateWallet()
	fmt.Printf("create address:%s",add)
	adds:=ws.GetAddresses()
	fmt.Printf("create addresses:%v",adds)
	ws.SaveToFile(wallet_file)
}

func TestWallets_LoadWalletsFromFile(t *testing.T) {
	ws, err := NewWallets(wallet_file)
	if err != nil {
		t.Error(err)
	}
	adds := ws.GetAddresses()
	fmt.Printf("create addresses:%v\n", adds)
}