package cli

import (
	"goclisandbox/cli/clog"
	"goclisandbox/cli/cmd"
	"goclisandbox/cli/git"

	"github.com/spf13/cobra"
)

var checkUpdateCmd = &cobra.Command{
	Use:     "u",
	Aliases: []string{"update"},
	Short:   "Update the given path 'f' must be given",
	Long:    "Check if the given path and tell if there is an update available",
	Run: func(cobraCommande *cobra.Command, args []string) {
		hasUpdate, err := isUpdateAvailable(args)

		if err != nil {
			panic("An error occured when checking if an update is available")
		}

		if hasUpdate {
			clog.ILog("An update is available")
		} else {
			clog.SLog("No update available")
		}
	},
}

func init() {
	rootCmd.AddCommand(checkUpdateCmd)
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
