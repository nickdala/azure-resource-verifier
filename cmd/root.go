/*
Copyright © 2024 Nick Dalalelis
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azure-resource-verifier",
	Short: "azure-resource-verifier is a tool that verifies Azure resources can be deployed to a region",
	Long: `azure-resource-verifier is a tool that verifies Azure resources can be deployed to a region. 
For more information, please visit https://github.com/nickdala/azure-resource-verifier`,

	SilenceErrors: true, // don't print errors twice, we handle them in cli.ExitOnError

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SetErr(os.Stderr)
	rootCmd.SetOut(os.Stdout)

	err := rootCmd.Execute()
	cli.ExitOnError(err, "")
}

func init() {
	//cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/.azure-resource-verifier.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		//home, err := os.UserHomeDir()
		//cobra.CheckErr(err)

		// Search config in home directory with name ".azure-resource-verifier" (without extension).
		//viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".azure-resource-verifier")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
