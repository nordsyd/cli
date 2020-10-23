package cmd

import (
	"fmt"

	"github.com/nordsyd/cli/core/manifest"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create a deploy to linked site",
	Long:  `This command will create and publish a deploy to the linked site`,
	Run: func(cmd *cobra.Command, args []string) {
		var files, err = manifest.GetFilesFromDir(args[0])

		if err != nil {
			panic(err)
		}

		fmt.Print(files)
	},
}
