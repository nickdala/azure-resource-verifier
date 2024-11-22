/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/nickdala/azure-resource-verifier/internal/cli/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// redisCmd represents the redis command
var redisCmd = &cobra.Command{
	Use:   "redis",
	Short: "Verify Azure Cache for Redis capabilities",
	Long:  `The redis command provides the means to verify Azure Cache for Redis capabilities.`,

	RunE: cli.AzureClientWrapRunE(redisCommand),
}

func redisCommand(cmd *cobra.Command, _ []string, cred *azidentity.DefaultAzureCredential, ctx context.Context) error {
	fmt.Println("redis called")

	subscriptionId := viper.GetString("subscription-id")
	log.Printf("subscription-id: %s", subscriptionId)

	clientFactory, err := armresources.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return cli.CreateAzrErr("failed to create the arm resource client factory", err)
	}

	res, err := clientFactory.NewProvidersClient().Get(ctx, "Microsoft.Cache", &armresources.ProvidersClientGetOptions{Expand: nil})
	if err != nil {
		return cli.CreateAzrErr("failed to get the cache provider", err)
	}

	redisLocations, err := getRedisLocations(&res.Provider)
	if err != nil {
		return cli.CreateAzrErr("failed to get the Azure Cache for Redis locations", err)
	}

	azureLocations, err := getLocations(cmd, cred, ctx, subscriptionId)
	if err != nil {
		return cli.CreateAzrErr("Error parsing location flag", err)
	}

	displayNameToLocation := make(map[string]string)
	for _, location := range azureLocations {
		displayNameToLocation[location.DisplayName] = location.Name
	}

	seenRegions := make(map[string]struct{})

	var data [][]string

	for _, location := range redisLocations {
		if regionName, ok := displayNameToLocation[*location]; ok {
			seenRegions[*location] = struct{}{}
			data = append(data, []string{regionName, *location, "true"})
		}
	}

	// Now add the regions that were not returned by the API
	for regionDisplayName, regionName := range displayNameToLocation {
		if _, ok := seenRegions[regionDisplayName]; !ok {
			data = append(data, []string{regionName, regionDisplayName, "false"})
		}
	}

	table := util.NewTable(util.RedisService)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func getRedisLocations(provider *armresources.Provider) ([]*string, error) {
	if provider.ResourceTypes == nil {
		return nil, cli.CreateAzrErr("failed to get the cache provider resource types", nil)
	}

	for _, resourceType := range provider.ResourceTypes {
		if resourceType.ResourceType == nil {
			//log.Printf("Skipping resource type with nil ResourceType")
			continue
		}

		// We're looking for locations for Redis
		if *resourceType.ResourceType != "Redis" {
			continue
		}

		if resourceType.Locations == nil {
			continue
		}

		return resourceType.Locations, nil
	}

	return nil, cli.CreateAzrErr("No Redis locations found", nil)
}

func init() {
	rootCmd.AddCommand(redisCmd)

	redisCmd.Flags().StringP("subscription-id", "s", "", "The Azure subscription id")
	// Required
	if err := redisCmd.MarkFlagRequired("subscription-id"); err != nil {
		redisCmd.Printf("Error marking flag required: %s", err)
		os.Exit(1)
	}

	redisCmd.Flags().StringArrayP("location", "l", []string{}, "The Azure location to list the capabilities. Can be specified multiple times")
	redisCmd.Flags().Bool("all-locations", false, "Whether to list capabilities for all locations")
	redisCmd.MarkFlagsOneRequired("location", "all-locations")
	redisCmd.MarkFlagsMutuallyExclusive("location", "all-locations")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickstartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickstartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
