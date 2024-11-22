/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers"
	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/nickdala/azure-resource-verifier/internal/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// postgresqlCmd represents the postgresql command
var postgresqlCmd = &cobra.Command{
	Use:   "postgresql",
	Short: "Verify Azure PostgreSQL Flexible Server capabilities",
	Long:  `The postgresql command provides the means to verify Azure PostgreSQL Flexible Server capabilities.`,

	RunE: cli.AzureClientWrapRunE(postgresqlCommand),
}

func postgresqlCommand(cmd *cobra.Command, _ []string, cred *azidentity.DefaultAzureCredential, ctx context.Context) error {
	fmt.Println("postgresql called")

	subscriptionId := viper.GetString("subscription-id")
	log.Printf("subscription-id: %s", subscriptionId)

	client, err := armpostgresqlflexibleservers.NewLocationBasedCapabilitiesClient(subscriptionId, cred, nil)
	if err != nil {
		return cli.CreateAzrErr("failed to create the postgresql flexible server client", err)
	}

	locations, err := getLocations(cmd, cred, ctx, subscriptionId)
	if err != nil {
		return cli.CreateAzrErr("Error parsing location flag", err)
	}

	var data [][]string

	for _, location := range locations {
		pager := client.NewExecutePager(location, nil)
		log.Printf("Getting capabilities for location %s", location)
		for pager.More() {
			nextResult, err := pager.NextPage(ctx)
			if err != nil {
				if azureErr, ok := err.(*azcore.ResponseError); ok {
					data = append(data, []string{location, "false", "false", azureErr.ErrorCode})
				} else {
					data = append(data, []string{location, "false", "false", err.Error()})
				}
				break
			}

			if len(nextResult.Value) == 0 {
				data = append(data, []string{location, "false", "false", "can't deploy to this location"})
				break
			}

			for _, capability := range nextResult.Value {
				data = append(data, []string{location, "true", strconv.FormatBool(*capability.ZoneRedundantHaSupported), ""})
				break // Only print the first capability
			}
		}
	}

	table := util.NewTable(util.PostgreSqlService)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func init() {
	rootCmd.AddCommand(postgresqlCmd)

	postgresqlCmd.Flags().StringP("subscription-id", "s", "", "The Azure subscription id")
	// Required
	if err := postgresqlCmd.MarkFlagRequired("subscription-id"); err != nil {
		postgresqlCmd.Printf("Error marking flag required: %s", err)
		os.Exit(1)
	}

	postgresqlCmd.Flags().StringArrayP("location", "l", []string{}, "The Azure location to list the capabilities. Can be specified multiple times")
	postgresqlCmd.Flags().Bool("all-locations", false, "Whether to list capabilities for all locations")
	postgresqlCmd.MarkFlagsOneRequired("location", "all-locations")
	postgresqlCmd.MarkFlagsMutuallyExclusive("location", "all-locations")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickstartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickstartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}