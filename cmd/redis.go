package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/nickdala/azure-resource-verifier/internal/azure"
	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/nickdala/azure-resource-verifier/internal/table"
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

	azureLocations, err := getLocations(cmd, cred, ctx, subscriptionId)
	if err != nil {
		return cli.CreateAzrErr("Error parsing location flag", err)
	}

	redisCache := azure.NewAzureRedisCache(cred, ctx, subscriptionId)
	redisLocations, err := redisCache.GetRedisLocations()
	if err != nil {
		return cli.CreateAzrErr("Error getting Redis locations", err)
	}

	deployableLocations := azureLocations.Intersection(redisLocations)
	unsupportedRegions := azureLocations.Difference(redisLocations)

	var data [][]string

	for _, location := range deployableLocations.Value {
		data = append(data, []string{location.Name, location.DisplayName, "true"})
	}

	for _, location := range unsupportedRegions.Value {
		data = append(data, []string{location.Name, location.DisplayName, "false"})
	}

	table := table.NewTable(table.RedisService)
	table.AppendBulk(data)
	table.Render()

	return nil
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
