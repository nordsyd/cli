package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/nordsyd/cli/core/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(authCmd)
}

var authQuestions = []*survey.Question{
	{
		Name:   "email",
		Prompt: &survey.Input{Message: "What is your email address?"},
		Validate: func(val interface{}) error {
			re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

			str := fmt.Sprintf("%v", val)

			if !re.MatchString(str) {
				return errors.New("Invalid email")
			}

			return nil
		},
		Transform: survey.ToLower,
	},
	{
		Name:     "password",
		Prompt:   &survey.Password{Message: "What is your password?"},
		Validate: survey.Required,
	},
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate user",
	Long:  `This command will fetch a JWT for future requests`,
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure there is no JWT on record
		if viper.Get("JWT") != "" {
			color.Yellow("You are already authenticated!")
			return
		}

		whiteBold := color.New(color.FgWhite).Add(color.Bold)
		whiteBold.Println("Sign in using your Nordsyd account")
		whiteBold.Println("Press Ctrl+C to cancel.")

		fmt.Print("\n")

		credentials := struct {
			Email    string
			Password string
		}{}

		err := survey.Ask(authQuestions, &credentials)
		if err != nil {
			if err == terminal.InterruptErr {
				color.Red("\nCancelling authentication")
				return
			}

			panic(err)
		}

		fmt.Println()

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)       // Build our new spinner
		s.Start()                                                         // Start the spinner
		jwt, error := api.GetJWT(credentials.Email, credentials.Password) // Run for some time to simulate work
		s.Stop()

		if error != nil {
			red := color.New(color.FgRed)
			red.Print("x ")
			whiteBold.Print(error)

			fmt.Print("\n")

			return
		}

		viper.Set("JWT", jwt)
		viper.WriteConfig()

		green := color.New(color.FgGreen)
		green.Print("✓ ")
		whiteBold.Print("Authentication successful!\n\n")

		fmt.Print("Next up, go to a folder with a static site, and type ")
		whiteBold.Print("nordsyd link")
		fmt.Print(" to link a site.\n")
	},
}
