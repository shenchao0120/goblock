package client

import (
	"github.com/spf13/cobra"
	"errors"
	"chaoshen.com/goblock/address"
	"strings"
)


func WalletCmd() *cobra.Command{
	//pflag:=sendCommand.PersistentFlags()
	walletCommand.AddCommand(walletCreateCommand)
	walletCommand.AddCommand(walletListCommand)


	return walletCommand
}

var walletCommand = &cobra.Command{
	Use:"wallet",
	Short:"manage wallet",
	Long:"manage wallet",
}


var walletCreateCommand = &cobra.Command{
	Use:"create",
	Short:"create new address",
	Long:"create new address",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("Wrong args number.")
		}
		ws,err:=address.NewWallets("")
		if err != nil {
			return err
		}
		addr:=ws.CreateWallet()
		ws.SaveToFile("")
		logger.Info("You create new address:%s",addr)
		return nil
	},
}

var walletListCommand = &cobra.Command{
	Use:"list",
	Short:"list all address",
	Long:"list all address",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New("Wrong args number.")
		}
		ws,err:=address.NewWallets("")
		if err != nil {
			return err
		}
		addrs:=ws.GetAddresses()
		logger.Infof("All address:\n%s",strings.Join(addrs,"\n"))
		return nil
	},
}



