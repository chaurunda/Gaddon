package cli

import (
	"fmt"
	"goclisandbox/cli/cmd"
	"goclisandbox/cli/git"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "u",
	Aliases: []string{"update"},
	Short:   "Update the given path",
	Long:    "Check if the given path and tell if there is an update available",
	Args:    cobra.ExactArgs(1),
	Run: func(cobraCommande *cobra.Command, args []string) {
		hasUpdate, err := isUpdateAvailable(args)

		if err != nil {
			panic("An error occured when checking if an update is available")
		}

		if hasUpdate {
			fmt.Println("An update is available")
		} else {
			fmt.Println("No update available")
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func isUpdateAvailable(args []string) (bool, error) {
	entries, err := cmd.ListDirectoryContents()
	isGitDir := false
	if err != nil {
		return false, err
	}
	for _, entry := range entries {
		if entry.IsGitDir() {
			isGitDir = true
		}
	}
	if isGitDir {
		lastLocalId := git.Check_last_local_commit_id(args[0])
		lastDistantId := git.Check_last_distant_commit_id(args[0])

		return !(lastDistantId == lastLocalId), nil
	}
	return false, nil
}
