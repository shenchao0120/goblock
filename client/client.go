package client

import (
	"github.com/spf13/cobra"
	"errors"
	"chaoshen.com/goblock/core"
	"github.com/op/go-logging"
)

var logger = logging.MustGetLogger("client")

func Cmd() *cobra.Command {
	chainCommand.AddCommand(SendCmd())
	chainCommand.AddCommand(getbalanceCmd())
	chainCommand.AddCommand(getInfoCmd())
	chainCommand.AddCommand(PrintCommand)
	chainCommand.AddCommand(CreateCommand)

	return chainCommand
}


var chainCommand = &cobra.Command{
	Use:"chain",
	Short:"Chain operation: creat |send | getbalance |print",
	Long:"Chain operation: creat |send | getbalance |print",
}




var CreateCommand = &cobra.Command{
	Use:"create",
	Short:"Create a New blockchain",
	Long:"Create a New blockchain",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Wrong args number.")
		}
		bc:=core.CreateBlockchain(args[0])
		if bc == nil {
			return errors.New("Create block chain error.")
		}
		logger.Info("Create new block chain success")
		return nil
	},
}

var PrintCommand = &cobra.Command{
	Use:"print",
	Short:"Print the chain ",
	Long:"Print the chain ",
	RunE: func(cmd *cobra.Command, args []string) error{
		if len(args)>0 {
			return errors.New("Wrong args number.")
		}
		bc:=core.NewBlockchain("")
		bc.Print()
		return nil
	},
}

//getbalance
//createblockchain
//send
//printchain
//address

