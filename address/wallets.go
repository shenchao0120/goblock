package address

import (
	"os"
	"io/ioutil"
	"github.com/op/go-logging"
	"encoding/gob"
	"crypto/elliptic"
	"bytes"
	"log"
)

var logger= logging.MustGetLogger("WALLETS")


const wallet_file = "wallet.dat"

type Wallets struct {
	ws map[string]*Wallet
}

func NewWallets(walletFile string) (*Wallets,error){
	if walletFile=="" {
		walletFile=wallet_file
	}
	wallets := &Wallets{}
	wallets.ws = make(map[string]*Wallet)
	err:=wallets.LoadWalletsFromFile(walletFile)
	if err != nil {
		return nil,err
	}
	return wallets,nil
}

func (ws *Wallets) LoadWalletsFromFile(walletFile string) error {
	if walletFile=="" {
		walletFile=wallet_file
	}
	if _,err:=os.Stat(walletFile); err != nil{
		if os.IsNotExist(err){
			logger.Info("The Wallets file not exist.")
			return nil
		}else{
			return err
		}
	}
	fileContent,err:=ioutil.ReadFile(walletFile)
	if err != nil {
		return err
	}
	gob.Register(elliptic.P256())
	decoder:=gob.NewDecoder(bytes.NewReader(fileContent))
	err=decoder.Decode(&ws.ws)
	if err != nil{
		return err
	}
	return nil
}

func (ws *Wallets) CreateWallet() string {
	wt:=NewWallet()
	address:=wt.GetAddress()
	logger.Info("Create wallet address:",address)
	ws.ws[address]=wt
	return address
}

func (ws *Wallets) GetWallet (address string) *Wallet {
	return ws.ws[address]
}

func (ws *Wallets) GetAddresses () []string{
	var addresses []string
	for add,_:=range ws.ws {
		addresses= append(addresses,add)
	}
	return addresses
}


func (ws *Wallets)SaveToFile(walletFile string){
	if walletFile=="" {
		walletFile=wallet_file
	}
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder:=gob.NewEncoder(&content)
	err:=encoder.Encode(ws.ws)
	if err != nil {
		log.Panic(err)
	}
	err=ioutil.WriteFile(walletFile,content.Bytes(),0644)
	if err != nil {
		log.Panic(err)
	}
}

