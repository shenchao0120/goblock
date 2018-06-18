package main

import (
	"chaoshen.com/goblock/core"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"os"
)

var logger = logging.MustGetLogger("New block.")

const cmdRoot = "core"
var versionFlag bool

var mainCmd = &cobra.Command{
	Use: "goblock",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

func main() {
	mainFlags:=mainCmd.PersistentFlags()
	mainFlags.BoolVarP(&versionFlag,"version","v",false,"Display current version of server.")

	mainCmd.AddCommand(core.Cmd())
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
	logger.Info("Exiting.....")
	/*
	bc := core.NewBlockchain()
	defer bc.Close()
	err := bc.AddBlock([]byte("first message"))
	CheckErr(err)
	err=bc.AddBlock([]byte("second message"))
	CheckErr(err)
	err=bc.AddBlock([]byte("third message"))
	CheckErr(err)
	bc.Print()
	*/
}

func CheckErr(err error) {
	if err != nil {
		logger.Error(err)
	}
	return
}
