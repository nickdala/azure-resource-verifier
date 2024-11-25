package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/nickdala/azure-resource-verifier/internal/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// locationlistCmd represents the locationlist command
var locationlistCmd = &cobra.Command{
	Use:   "list-locations",
	Short: "Lists all locations in the Azure subscription",
	Long:  `The list-locations command lists all locations in the Azure subscription.`,
	RunE:  cli.AzureClientWrapRunE(listLocationsCommand),
}

func listLocationsCommand(cmd *cobra.Command, _ []string, cred *azidentity.DefaultAzureCredential, ctx context.Context) error {
	fmt.Println("Listing all locations in the Azure subscription")

	subscriptionId := viper.GetString("subscription-id")
	log.Printf("subscription-id: %s", subscriptionId)

	locations, err := getAllLocationsFromSubscription(cred, ctx, subscriptionId)
	if err != nil {
		return cli.CreateAzrErr("Error getting locations", err)
	}

	table := table.NewTable(table.Locations)

	for _, location := range locations.Value {
		table.AppendRow([]string{location.Name, location.DisplayName})
	}

	table.Render()

	return nil

}

func init() {
	rootCmd.AddCommand(locationlistCmd)

	locationlistCmd.Flags().StringP("subscription-id", "s", "", "The Azure subscription id")
	// Required
	if err := locationlistCmd.MarkFlagRequired("subscription-id"); err != nil {
		locationlistCmd.Printf("Error marking flag required: %s", err)
		os.Exit(1)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// locationlist.goCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// locationlist.goCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
