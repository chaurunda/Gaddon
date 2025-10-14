package cli

import (
	"goclisandbox/cli/clog"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:     "i",
	Aliases: []string{"install"},
	Short:   "Install new addon",
	Long:    "Install and rename the folder to be the same as the .toc file inside it",
	Args:    cobra.ExactArgs(1), //folder path and install
	Run: func(cmd *cobra.Command, args []string) {
		clog.Log("Install in progress")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
