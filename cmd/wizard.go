package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers"
	"github.com/nickdala/azure-resource-verifier/cmd/modal/appservice"
	"github.com/nickdala/azure-resource-verifier/cmd/modal/database"
	azrazure "github.com/nickdala/azure-resource-verifier/internal/azure"
	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// wizardCmd represents the wizard command
var wizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Wizard to guide you through the Azure Resource Verifier",
	Long:  `The wizard command provides a guided experience that shows the regions the Azure resources can be deployed.`,
	RunE:  cli.AzureClientWrapRunE(wizardCommand),
}

func wizardCommand(cmd *cobra.Command, _ []string, cred *azidentity.DefaultAzureCredential, ctx context.Context) error {
	fmt.Println("Listing the locations where the Azure resources can be deployed")

	subscriptionId := viper.GetString("subscription-id")
	log.Printf("subscription-id: %s", subscriptionId)

	/*azureLocations, err := getLocations(cmd, cred, ctx, subscriptionId)
	if err != nil {
		return cli.CreateAzrErr("Error parsing location flag", err)
	}*/

	//var postgresLocations []*azrazure.AzureLocation
	//var redisLocations []*azrazure.AzureLocation

	databases, _ := database.ShowDatabaseModalAndGetChoices()
	for _, db := range databases {
		switch db {
		case database.REDIS:
			println("Selected: Azure Cache for Redis")
		case database.POSTGRESQL:
			println("Selected: Azure PostgreSQL Flexible Server")
			/*postgresLocations, err = getPostgresLocations(subscriptionId, cred, ctx, azureLocations, false)
			if err != nil {
				return cli.CreateAzrErr("Error getting PostgreSQL locations", err)
			}*/
		case database.POSTGRESQL_HA:
			println("Selected: Azure PostgreSQL Flexible Server with HA")
			/*postgresLocations, err = getPostgresLocations(subscriptionId, cred, ctx, azureLocations, true)
			if err != nil {
				return cli.CreateAzrErr("Error getting PostgreSQL with HA locations", err)
			}*/
		}
	}

	appService, _ := appservice.ShowAppServiceModalAndGetChoices()
	switch appService {
	case appservice.APP_SERVICE_LINUX_CODE:
		println("Selected: Azure App Service - Linux Code")
	case appservice.APP_SERVICE_LINUX_CONTAINER:
		println("Selected: Azure App Service - Linux Container")
	case appservice.APP_SERVICE_WINDOWS_CODE:
		println("Selected: Azure App Service - Windows Code")
	case appservice.APP_SERVICE_WINDOWS_CONTAINER:
		println("Selected: Azure App Service - Windows Container")
	}

	//deployableLocations := make([]*azrazure.AzureLocation, 0)

	return nil

}

func getPostgresLocations(subscriptionId string, cred *azidentity.DefaultAzureCredential, ctx context.Context, filter []*azrazure.AzureLocation, haEnabled bool) ([]*azrazure.AzureLocation, error) {
	locations := make([]*azrazure.AzureLocation, 0)

	client, err := armpostgresqlflexibleservers.NewLocationBasedCapabilitiesClient(subscriptionId, cred, nil)
	if err != nil {
		return locations, cli.CreateAzrErr("failed to create the postgresql flexible server client", err)
	}

	for _, location := range filter {
		pager := client.NewExecutePager(location.Name, nil)
		for pager.More() {
			nextResult, err := pager.NextPage(ctx)
			if err != nil {
				break
			}

			if len(nextResult.Value) == 0 {
				break
			}

			// We have the capabilities for the location.
			// You can at least deploy PostgreSQL Flexible Server to this location.
			// If HA is enabled, check if the location supports HA.
			if !haEnabled {
				locations = append(locations, location)
				break
			}

			// Check if the location supports HA
			for _, capability := range nextResult.Value {
				if *capability.ZoneRedundantHaSupported {
					locations = append(locations, location)
					break // Only need confirmation for one capability for HA
				}
			}
		}
	}

	return locations, nil
}

func getRedisLocations2(subscriptionId string, cred *azidentity.DefaultAzureCredential, ctx context.Context, filter []*azrazure.AzureLocation) ([]*azrazure.AzureLocation, error) {
	locations := make([]*azrazure.AzureLocation, 0)

	return locations, nil
}

func init() {
	rootCmd.AddCommand(wizardCmd)

	wizardCmd.Flags().StringP("subscription-id", "s", "", "The Azure subscription id")
	// Required
	if err := wizardCmd.MarkFlagRequired("subscription-id"); err != nil {
		wizardCmd.Printf("Error marking flag required: %s", err)
		os.Exit(1)
	}

	wizardCmd.Flags().StringArrayP("location", "l", []string{}, "The Azure location to list the capabilities. Can be specified multiple times")
	wizardCmd.Flags().Bool("all-locations", false, "Whether to list capabilities for all locations")
	wizardCmd.MarkFlagsOneRequired("location", "all-locations")
	wizardCmd.MarkFlagsMutuallyExclusive("location", "all-locations")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wizard.goCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wizard.goCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
