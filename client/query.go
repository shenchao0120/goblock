package client

import (
	"github.com/spf13/cobra"
	"chaoshen.com/goblock/core"
	"errors"
)

func getbalanceCmd() *cobra.Command{
	return balanceCmd
}


var balanceCmd = &cobra.Command{
	Use:"getbalance",
	Short:"get the balance of specific address",
	Long:"get the balance of specific address",
	RunE:getBalance,
}


func getBalance(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("Wrong args number.")
	}
	bc:=core.GetBlockchain()
	balance:=bc.GetBalance(args[0])
	logger.Infof("The balance of address[%s] is [%d]",args[0],balance)
	return nil
}

func getInfoCmd() *cobra.Command{
	return infocmd
}

var infocmd = &cobra.Command{
	Use:"getinfo",
	Short:"get the info of the blockchain",
	Long:"get the info of the blockchai",
	RunE: func(cmd *cobra.Command, args []string) error {
		bc:=core.GetBlockchain()
		logger.Infof("Genesis Block Data:%s",string(bc.GetBlock(1).Transactions[0].Vin[0].PubKey))
		logger.Infof("Height [%d]",bc.Height)
		logger.Infof("Transaction Num [%d]",bc.GetTransactionNum())
		return nil
	},
}
