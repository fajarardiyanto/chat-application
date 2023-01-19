package main

import (
	"github.com/fajarardiyanto/chat-application/cmd/api"
	"github.com/fajarardiyanto/chat-application/cmd/generated"
	"github.com/fajarardiyanto/chat-application/config"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	if err := Run(os.Args[1:]); err != nil {
		config.GetLogger().Error("Unable to run the command %s ", err.Error())
	}
}

var rootCmd = &cobra.Command{
	Use:   config.GetConfig().Name,
	Short: config.GetConfig().Name,
}

func init() {
	rootCmd.AddCommand(api.CmdAPI)
	rootCmd.AddCommand(generated.CmdConfig)
}

func Run(args []string) error {
	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}
