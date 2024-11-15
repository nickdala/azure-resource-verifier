/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// quickstartCmd represents the quickstart command
var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Wizard to quickly get started with azure-resource-verifier",
	Long: `The quickstart command provide the means to quickly get started with azure-resource-verifier.
It will guide you through the process of creating a new configuration file and running a verification.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("quickstart called")
	},
}

func init() {
	rootCmd.AddCommand(quickstartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickstartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickstartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
