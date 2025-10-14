package cli

import (
	"goclisandbox/cli/clog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		clog.ELog("To run Gaddon you have to give a flag and a path ex : gaddon . -u /path/to/folder")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		clog.ELog("Oops. An error while executing Zero: " + err.Error())
		os.Exit(1)
	}
}
