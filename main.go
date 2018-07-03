package main

import (
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
	"os"
	"chaoshen.com/goblock/client"
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

	mainCmd.AddCommand(client.Cmd())
	mainCmd.AddCommand(client.WalletCmd())
	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
	logger.Info("Exiting.....")
}

func CheckErr(err error) {
	if err != nil {
		logger.Error(err)
	}
	return
}
