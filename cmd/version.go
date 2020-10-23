package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is version
var Version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the Nordsyd CLI",
	Long:  `Print the Nordsyd CLI version number for future reference.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nordsyd CLI", Version)
	},
}
