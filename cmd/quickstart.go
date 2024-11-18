/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers"
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
	fmt.Println("quickstart called **")

	subscriptionId := viper.GetString("subscription-id")
	log.Printf("subscription-id: %s", subscriptionId)

	client, err := armpostgresqlflexibleservers.NewLocationBasedCapabilitiesClient(subscriptionId, cred, nil)
	if err != nil {
		return cli.CreateAzrErr("failed to create client", err)
	}

	locations, err := cmd.Flags().GetStringArray("location")
	if err != nil {
		return cli.CreateAzrErr("Error parsing location flag", err)
	}

	for _, location := range locations {
		pager := client.NewExecutePager(location, nil)
		log.Printf("Listing capabilities for location %s", location)
		for pager.More() {
			nextResult, err := pager.NextPage(ctx)
			if err != nil {
				return cli.CreateAzrErr("failed to get next page", err)
			}

			if len(nextResult.Value) == 0 {
				log.Printf("no capabilities found")
				break
			}

			log.Println("Capabilities:")
			for _, capability := range nextResult.Value {
				log.Printf("Zone: %v Status: %v HA: %v\n", *capability.Zone, *capability.Status, *capability.ZoneRedundantHaSupported)
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(quickstartCmd)

	quickstartCmd.Flags().StringP("subscription-id", "s", "", "The Azure subscription id")
	quickstartCmd.Flags().StringArrayP("location", "l", []string{}, "The Azure location. Can be specified multiple times")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickstartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickstartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
