package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// quickstartCmd represents the quickstart command
var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Quickstart azure-resource-verifier",
	Long: `The quickstart command provide the means to quickly get started with azure-resource-verifier.
It will guide you through the process of creating a new configuration file and running a verification.`,

	RunE: cli.AzureClientWrapRunE(quickStartCommand),
}

func quickStartCommand(cmd *cobra.Command, _ []string, cred *azidentity.DefaultAzureCredential, ctx context.Context) error {
	fmt.Println("quickstart called")

	subscriptionId := viper.GetString("subscription-id")
	log.Printf("subscription-id: %s", subscriptionId)

	return nil
}

func init() {
	rootCmd.AddCommand(quickstartCmd)

	quickstartCmd.Flags().StringP("subscription-id", "s", "", "The Azure subscription id")
	// Required
	if err := quickstartCmd.MarkFlagRequired("subscription-id"); err != nil {
		quickstartCmd.Printf("Error marking flag required: %s", err)
		os.Exit(1)
	}

	quickstartCmd.Flags().StringArrayP("location", "l", []string{}, "The Azure location to list the capabilities. Can be specified multiple times")
	quickstartCmd.Flags().Bool("all-locations", false, "Whether to list capabilities for all locations")
	quickstartCmd.MarkFlagsOneRequired("location", "all-locations")
	quickstartCmd.MarkFlagsMutuallyExclusive("location", "all-locations")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickstartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickstartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
