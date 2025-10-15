package cli

import (
	"bufio"
	"goclisandbox/cli/clog"
	"goclisandbox/cli/cmd"
	"goclisandbox/cli/git"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var checkUpdateCmd = &cobra.Command{
	Use:     "u",
	Aliases: []string{"update"},
	Short:   "Update the given path 'f' must be given",
	Long:    "Check if the given path and tell if there is an update available",
	Run: func(cobraCommande *cobra.Command, args []string) {
		filepath, err := cobraCommande.Flags().GetString("folder")
		if err != nil {
			clog.ELog("Error reading folder flag")
			return
		}

		// If no folder flag provided, prompt for it
		if filepath == "" {
			filepath = PromptForPath("Please give a wow addon folder path:")
		}

		hasUpdate, err := isUpdateAvailable(filepath)

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
	checkUpdateCmd.Flags().String("folder", "", "path to the wow addon folder")
	rootCmd.AddCommand(checkUpdateCmd)
}

func isUpdateAvailable(path string) (bool, error) {
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
		lastLocalId := git.Check_last_local_commit_id(path)
		lastDistantId := git.Check_last_distant_commit_id(path)

		return !(lastDistantId == lastLocalId), nil
	}
	return false, nil
}

func PromptForPath(promptMessage string) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		clog.ILog(promptMessage)
		input, err := reader.ReadString('\n')
		if err != nil {
			clog.ELog("Error reading input")
			continue
		}

		// Remove newline and whitespace
		path := strings.TrimSpace(input)

		// Check if path is not empty
		if path == "" {
			clog.WLog("Path cannot be empty, please try again")
			continue
		}

		// Validate the path exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			clog.WLog("Path does not exist, please try again")
			continue
		}

		// Validate it's a directory
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			clog.WLog("Path must be a directory, please try again")
			continue
		}

		clog.SLog("Valid path provided: " + path)
		return path
	}
}
