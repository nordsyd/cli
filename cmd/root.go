package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "nordsyd",
	Short: "Nordsyd CLI used for interacting with the Nordsyd API",
	Long:  `With the Nordsyd CLI you can deploy your static site to our global CDN network`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func getConfigDirectory() string {
	homeDir, error := homedir.Dir()

	if error != nil {
		panic("Unable to determine home directory")
	}

	return homeDir + "/.nordsyd"
}

func configDirectoryExists() bool {
	folderStat, error := os.Stat(getConfigDirectory())

	if error != nil {
		return false
	}

	return folderStat.IsDir()
}

func init() {
	cobra.OnInitialize(initConfig)

	viper.SetDefault("JWT", "")
	viper.SetDefault("SiteLinks", map[string]string{})
}

func initConfig() {
	configDirectory := getConfigDirectory()

	if !configDirectoryExists() {
		err := os.Mkdir(configDirectory, 0755)

		if err != nil {
			log.Fatal(err)
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(configDirectory)

	if err := viper.ReadInConfig(); err != nil {
		//fmt.Println("Failed to read current config (it's probably missing)")

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SafeWriteConfig()
			//fmt.Println("Created new user config")

			//fmt.Println(err)
		} else {
			panic("Config file not found, but something else happened")
		}
	}

	//viper.Set("JWT", "")
	//
	//viper.WriteConfig()

	//fmt.Println(viper.AllSettings())
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
