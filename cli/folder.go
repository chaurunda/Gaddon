package cli

import (
	"goclisandbox/cli/clog"

	"github.com/spf13/cobra"
)

var folderCmd = &cobra.Command{
	Use:     "f",
	Aliases: []string{"folder"},
	Short:   "The path you want to Graddon work into",
	Long:    "The path you want to Graddon work into, for update or install",
	Args:    cobra.ExactArgs(1), // The path
	Run: func(cmd *cobra.Command, args []string) {
		clog.Log("The path is " + args[0])
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)
}
