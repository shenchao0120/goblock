package client

import (
	"github.com/spf13/cobra"
	"errors"
	"chaoshen.com/goblock/core"
	"chaoshen.com/goblock/transaction"
)

var (
	fromAddr string
	toAddr string
	amount int
)

func SendCmd() *cobra.Command{
	pflag:=sendCommand.PersistentFlags()
	pflag.StringVarP(&fromAddr,"from","f","","The address of sender.")
	pflag.StringVarP(&toAddr,"to","t","","The address of receiver.")
	pflag.IntVarP(&amount,"amt","a",0,"The transfer amount.")
	return sendCommand
}

var sendCommand = &cobra.Command{
	Use:"send",
	Short:"Send AMOUNT of coins from FROM address to To",
	Long:"Send AMOUNT of coins from FROM address to To",
	RunE: Send ,
}


func Send(cmd *cobra.Command, args []string) error {
	if len(args) != 0{
		return errors.New("Wrong args number.")
	}
	bc:=core.GetBlockchain()

	tx,err:=bc.NewUTXOTransaction(fromAddr,toAddr,amount)
	if err != nil {
		return err
	}
	err=bc.GenerateBlock([]*transaction.Transaction{tx})
	if err != nil {
		return err
	}
	return nil
}

