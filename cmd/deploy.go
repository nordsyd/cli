package cmd

import (
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	"github.com/nordsyd/cli/core/api"
	"github.com/nordsyd/cli/core/cli"
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
		if !cli.EnsureAuthentication() {
			return
		}

		deployFolder := getCurrentDirectory()

		// Check if folder has been supplied
		if len(args) == 1 {
			deployFolder = path.Join(deployFolder, args[0])
		}

		whiteBold := color.New(color.FgWhite).Add(color.Bold)

		fmt.Print("Deploying folder: ")
		whiteBold.Print(deployFolder)
		fmt.Print("\nPress Ctrl+C to cancel.\n\n")

		// Get files from directory
		scanSpinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		scanSpinner.UpdateCharSet(spinner.CharSets[14])
		scanSpinner.Suffix = " Scanning folder for files"
		scanSpinner.Start()
		files, scanError := manifest.GetFilesFromDir(deployFolder)
		scanSpinner.Stop()

		if scanError != nil {
			cli.PrintError(scanError.Error())
			return
		}

		cli.PrintSuccess("Folder scan complete")

		// Generate manifest
		manifestSpinner := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
		manifestSpinner.UpdateCharSet(spinner.CharSets[14])
		manifestSpinner.Suffix = " Generating deploy manifest"
		manifestSpinner.Start()
		manifest := manifest.GenerateManifest(files, deployFolder)
		manifestSpinner.Stop()

		cli.PrintSuccess("Deploy manifest generated")

		siteSlug := getSiteLink(getCurrentDirectory())
		deployResponse := api.CreateDeploy(siteSlug, manifest)

		if len(deployResponse.ToUpload) == 0 {
			fmt.Print("\n")
			cli.PrintSuccess("Deploy is already active, nothing to do")
			return
		}

		// Perform upload
		count := len(deployResponse.ToUpload)
		bar := uiprogress.AddBar(count).AppendCompleted().PrependElapsed()
		bar.PrependFunc(func(b *uiprogress.Bar) string {
			return fmt.Sprintf("  Uploading files... (%d/%d)", b.Current(), count)
		})

		uiprogress.Start()
		var wg sync.WaitGroup

		for _, uploadFile := range deployResponse.ToUpload {
			fullPath := path.Join(deployFolder, uploadFile.Key)

			wg.Add(1)

			go func(URL string) {
				defer wg.Done()

				api.UploadFile(fullPath, URL)

				bar.Incr()
			}(uploadFile.UploadURL)
		}

		time.Sleep(time.Second)
		wg.Wait()
		uiprogress.Stop()

		//fmt.Println(generatedManifest)

		api.FinaliseDeploy(siteSlug, deployResponse.DeployHash)

		fmt.Print("\n")
		cli.PrintSuccess("Deploy successful")
	},
}
