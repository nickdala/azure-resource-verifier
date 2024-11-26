package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/nickdala/azure-resource-verifier/cmd/modal/appservice"
	"github.com/nickdala/azure-resource-verifier/cmd/modal/database"
	"github.com/nickdala/azure-resource-verifier/internal/azure"
	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/nickdala/azure-resource-verifier/internal/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// quickstartCmd represents the quickstart command
var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Quickstart Azure Resource Verifier",
	Long:  `The quickstart command command provides a guided experience that shows the regions the Azure resources can be deployed.`,

	RunE: cli.AzureClientWrapRunE(quickStartCommand),
}

func quickStartCommand(cmd *cobra.Command, _ []string, cred *azidentity.DefaultAzureCredential, ctx context.Context) error {
	fmt.Println("quickstart called")

	subscriptionId := viper.GetString("subscription-id")
	log.Printf("subscription-id: %s", subscriptionId)

	azureLocations, err := getLocations(cmd, cred, ctx, subscriptionId)
	if err != nil {
		return cli.CreateAzrErr("Error parsing location flag", err)
	}

	databases, _ := database.ShowDatabaseModalAndGetChoices()
	appService, _ := appservice.ShowAppServiceModalAndGetChoices()

	switch appService {
	case appservice.APP_SERVICE_LINUX_CODE:
		println("Selected: Azure App Service - Linux Code")
		azureLocations, err = getLocationsForAppService(azureLocations, cred, ctx, subscriptionId, azure.Linux, azure.Code)
		if err != nil {
			return cli.CreateAzrErr("Error getting App Service locations", err)
		}
	case appservice.APP_SERVICE_LINUX_CONTAINER:
		println("Selected: Azure App Service - Linux Container")
		azureLocations, err = getLocationsForAppService(azureLocations, cred, ctx, subscriptionId, azure.Linux, azure.Container)
		if err != nil {
			return cli.CreateAzrErr("Error getting App Service locations", err)
		}
	case appservice.APP_SERVICE_WINDOWS_CODE:
		println("Selected: Azure App Service - Windows Code")
		azureLocations, err = getLocationsForAppService(azureLocations, cred, ctx, subscriptionId, azure.Windows, azure.Code)
		if err != nil {
			return cli.CreateAzrErr("Error getting App Service locations", err)
		}
	case appservice.APP_SERVICE_WINDOWS_CONTAINER:
		println("Selected: Azure App Service - Windows Container")
		azureLocations, err = getLocationsForAppService(azureLocations, cred, ctx, subscriptionId, azure.Windows, azure.Container)
		if err != nil {
			return cli.CreateAzrErr("Error getting App Service locations", err)
		}
	}

	for _, db := range databases {
		switch db {
		case database.REDIS:
			println("Selected: Azure Cache for Redis")
			azureLocations, err = getLocationsForRedis(azureLocations, cred, ctx, subscriptionId)
			if err != nil {
				return cli.CreateAzrErr("Error getting Redis locations", err)
			}
		case database.POSTGRESQL:
			println("Selected: Azure PostgreSQL Flexible Server")
			azureLocations, err = getPostgresLocations(subscriptionId, cred, ctx, azureLocations, false)
			if err != nil {
				return cli.CreateAzrErr("Error getting PostgreSQL locations", err)
			}
		case database.POSTGRESQL_HA:
			println("Selected: Azure PostgreSQL Flexible Server with HA")
			azureLocations, err = getPostgresLocations(subscriptionId, cred, ctx, azureLocations, true)
			if err != nil {
				return cli.CreateAzrErr("Error getting PostgreSQL HA locations", err)
			}
		}
	}

	var data [][]string
	for _, location := range azureLocations.Value {
		data = append(data, []string{location.Name, location.DisplayName})
	}

	table := table.NewTable(table.Locations)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func getLocationsForAppService(locations *azure.AzureLocationList, cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string, os azure.AppServiceOS, publishType azure.AppServicePublishType) (*azure.AzureLocationList, error) {
	azureAppService := azure.NewAzureAppService(cred, ctx, subscriptionId)
	appServiceLocations, err := azureAppService.GetAppServiceLocations(locations, os, publishType)
	if err != nil {
		return nil, fmt.Errorf("error getting App Service locations %w", err)
	}

	return locations.Intersection(appServiceLocations), nil
}

func getLocationsForRedis(locations *azure.AzureLocationList, cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) (*azure.AzureLocationList, error) {
	redisCache := azure.NewAzureRedisCache(cred, ctx, subscriptionId)
	redisLocations, err := redisCache.GetRedisLocations()
	if err != nil {
		return nil, fmt.Errorf("error getting Redis locations %w", err)
	}

	return locations.Intersection(redisLocations), nil
}

func getPostgresLocations(subscriptionId string, cred *azidentity.DefaultAzureCredential, ctx context.Context, locations *azure.AzureLocationList, haEnabled bool) (*azure.AzureLocationList, error) {

	azurePostgresql := azure.NewAzurePostgresqlFlexibleServer(cred, ctx, subscriptionId)

	postgresLocations, postgresqlHaLocations, _, err := azurePostgresql.GetPostgresqlLocations(locations)
	if err != nil {
		return nil, fmt.Errorf("error getting PostgreSQL locations %w", err)
	}

	if haEnabled {
		return locations.Intersection((*azure.AzureLocationList)(postgresqlHaLocations)), nil
	} else {
		return locations.Intersection((*azure.AzureLocationList)(postgresLocations)), nil
	}
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
