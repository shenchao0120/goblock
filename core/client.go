package core

import (
	"github.com/spf13/cobra"
	"errors"
)

func Cmd() *cobra.Command {
	chainCommand.AddCommand(addCommand)
	chainCommand.AddCommand(printCommand)
	return chainCommand
}


var chainCommand = &cobra.Command{
	Use:"chain",
	Short:"Chain operation: add | print",
	Long:"Chain operation: add | print",
}


var addCommand = &cobra.Command{
	Use:"add",
	Short:"Add a block by data",
	Long:"Add a block by data",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Wrong args number.")
		}
		bc:=NewBlockchain()
		err:=bc.AddBlock([]byte(args[0]))
		if err != nil {
			return err
		}
		return nil
	},
}

var printCommand = &cobra.Command{
	Use:"print",
	Short:"Print the chain ",
	Long:"Print the chain ",
	RunE: func(cmd *cobra.Command, args []string) error{
		if len(args)>0 {
			return errors.New("Wrong args number.")
		}
		bc:=NewBlockchain()
		bc.Print()
		return nil
	},
}

