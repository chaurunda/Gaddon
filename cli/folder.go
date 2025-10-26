package cli

import (
	"goclisandbox/cli/clog"
	"log"
	"os"

	"github.com/joho/godotenv"
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
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		fileName := os.Getenv("GADDON_FILE_NAME")
		homeDirPath, err := os.UserHomeDir()

		if err != nil {
			panic("An error occured on reading user home dir")
		}

		filePath := homeDirPath + "/" + fileName

		// Check if file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// File doesn't exist, create it and write args[0]
			file, err := os.Create(filePath)
			if err != nil {
				clog.Log("Error creating gaddon.txt: " + err.Error())
				return
			}
			defer file.Close()

			_, err = file.WriteString(args[0])
			if err != nil {
				clog.Log("Error writing to gaddon.txt: " + err.Error())
				return
			}

			clog.Log("Created gaddon.txt and wrote path: " + args[0])
		} else {
			clog.Log("gaddon.txt already exists")
		}
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)
}
