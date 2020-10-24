package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/nordsyd/cli/core/cli"

	"github.com/fatih/color"
	"github.com/spf13/viper"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/nordsyd/cli/core/api"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(linkCmd)
}

func getUserSitesStrings(sites []api.Site) []string {
	listOfSites := []string{}

	for _, s := range sites {
		listOfSites = append(listOfSites, s.Name)
	}

	return listOfSites
}

func setSiteLink(directory string, siteSlug string) {
	siteLinks := viper.GetStringMapString("SiteLinks")
	siteLinks[directory] = siteSlug
	viper.Set("SiteLinks", siteLinks)
	viper.WriteConfig()
}

func getCurrentDirectory() string {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	return path
}

func getSiteLink(directory string) string {
	siteLinks := viper.GetStringMapString("SiteLinks")

	link, ok := siteLinks[directory]

	if !ok {
		return ""
	}

	return link
}

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Links a local site to Nordsyd site",
	Long:  `Links a local site to Nordsyd site, such that the CLI deploys to the correct site`,
	Run: func(cmd *cobra.Command, args []string) {
		if !cli.EnsureAuthentication() {
			return
		}

		currentDirectory := getCurrentDirectory()
		siteLink := getSiteLink(currentDirectory)

		override := false

		white := color.New(color.FgWhite).Add(color.Bold)

		if siteLink != "" {
			yellow := color.New(color.FgYellow).Add(color.Bold)
			yellow.Print("WARNING! ")

			white.Print("This directory is already linked to: ", siteLink)

			fmt.Print("\n\n")

			prompt := &survey.Confirm{
				Message: "Would you like to override the current link?",
			}
			err := survey.AskOne(prompt, &override)

			if err != nil {
				if err == terminal.InterruptErr {
					color.Red("\nCancelling site link")
					return
				}

				panic(err)
			}

			fmt.Print("\n\n")
		}

		if !override {
			white.Println("\nNo changes were made.")
			return
		}

		fmt.Println("Linking a site, allows you to deploy to the Nordsyd CDN")
		fmt.Println("Press Ctrl+C to cancel at any time")
		fmt.Print("\n")

		userSites := api.GetUserSites()
		userSiteStrings := getUserSitesStrings(userSites)

		var linkQuestions = []*survey.Question{
			{
				Name: "SiteIndex",
				Prompt: &survey.Select{
					Message: "Please select the site you'd like to link:",
					Options: userSiteStrings,
				},
				Validate:  survey.Required,
				Transform: survey.ToLower,
			},
		}

		answers := struct {
			SiteIndex int
		}{}

		err := survey.Ask(linkQuestions, &answers)
		if err != nil {
			panic(err)
		}

		pickedSite := userSites[answers.SiteIndex]

		setSiteLink(currentDirectory, pickedSite.Slug)

		green := color.New(color.FgGreen)
		green.Print("✓ ")
		white.Print("Linked site successfully!\n\n")

		fmt.Print("Next up, to create a deploy, type ")
		white.Print("nordsyd deploy")
		fmt.Print(", optionally followed by a subfolder containing your compiled site.\n")
	},
}
