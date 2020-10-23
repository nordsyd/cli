package cmd

import (
	"github.com/fatih/color"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets and removes all user configuration",
	Long:  `Used to reset and remove all traces of the Nordsyd CLI. Be careful.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("JWT", "")
		viper.WriteConfig()

		whiteBold := color.New(color.FgWhite).Add(color.Bold)
		green := color.New(color.FgGreen)
		green.Print("✓ ")
		whiteBold.Print("User configuration reset\n")
	},
}
