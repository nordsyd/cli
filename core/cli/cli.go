package cli

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

// EnsureAuthentication is a CLI helper to make sure user is signed in
func EnsureAuthentication() bool {
	if viper.Get("JWT") == "" {
		white := color.New(color.FgWhite)
		whiteBold := color.New(color.FgWhite).Add(color.Bold)
		yellow := color.New(color.FgYellow).Add(color.Bold)

		yellow.Print("WARNING! ")

		white.Print("You need to sign in before linking a site.\n\nUse ")
		whiteBold.Print("nordsyd auth")
		white.Print(" to sign in.\n")

		return false
	}

	return true
}

// PrintError is a helper function to print out errors
func PrintError(errorMessage string) {
	red := color.New(color.FgRed)
	whiteBold := color.New(color.FgWhite).Add(color.Bold)

	red.Print("x ")
	whiteBold.Print(errorMessage)

	fmt.Print("\n")
}

// PrintSuccess is a helper function to print out success messages
func PrintSuccess(message string) {
	green := color.New(color.FgGreen)
	whiteBold := color.New(color.FgWhite).Add(color.Bold)

	green.Print("✓ ")
	whiteBold.Print(message)

	fmt.Print("\n")
}
