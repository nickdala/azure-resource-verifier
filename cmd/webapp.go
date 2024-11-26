/*
Copyright © 2024 Nick Dalalelis
*/
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

var (
	webAppOperatingSystemChoice = cli.CliChoice{
		Name:        "operating-system",
		Description: "The operating system of the web app (linux or windows)",
		Default:     "linux",
		Choices:     []string{"linux", "windows"},
	}

	publishType = cli.CliChoice{
		Name:        "publish-type",
		Description: "The publish type of the web app (code or container)",
		Default:     "code",
		Choices:     []string{"code", "container"},
	}
)

// webAppCmd represents the web-app command
var webAppCmd = &cobra.Command{
	Use:   "web-app",
	Short: "Verify Azure App Service Web App can be deployed to a location",
	Long:  `The web-app command provides the means to verify if Azure App Service Web App can be deploy to a location.`,

	RunE: cli.AzureClientWrapRunE(appServiceCommand),
}

func appServiceCommand(cmd *cobra.Command, _ []string, cred *azidentity.DefaultAzureCredential, ctx context.Context) error {
	fmt.Println("web-app called")

	subscriptionId := viper.GetString("subscription-id")
	log.Printf("subscription-id: %s", subscriptionId)

	os := viper.GetString(webAppOperatingSystemChoice.Name)
	if valid := webAppOperatingSystemChoice.IsValidChoice(os); !valid {
		return cli.CreateAzrErr(fmt.Sprintf("Invalid operating system choice: %s", os), nil)
	}

	publish := viper.GetString(publishType.Name)
	if valid := publishType.IsValidChoice(publish); !valid {
		return cli.CreateAzrErr(fmt.Sprintf("Invalid publish type choice: %s", publish), nil)
	}

	osType, err := azure.AppServiceOSFromString(os)
	if err != nil {
		return cli.CreateAzrErr("Error parsing operating system flag", err)
	}

	publishType, err := azure.AppServicePublishTypeFromString(publish)
	if err != nil {
		return cli.CreateAzrErr("Error parsing publish type flag", err)
	}

	azureLocations, err := getLocations(cmd, cred, ctx, subscriptionId)
	if err != nil {
		return cli.CreateAzrErr("Error parsing location flag", err)
	}

	azureAppService := azure.NewAzureAppService(cred, ctx, subscriptionId)
	appServiceLocations, err := azureAppService.GetAppServiceLocations(azureLocations, osType, publishType)
	if err != nil {
		return cli.CreateAzrErr("Error getting App Service locations", err)
	}

	var data [][]string

	seenRegions := make(map[string]struct{})

	for _, location := range appServiceLocations.Value {
		data = append(data, []string{location.Name, location.DisplayName, "true"})
		seenRegions[location.Name] = struct{}{}
	}

	// Now add the regions that were not returned by the API
	for _, location := range azureLocations.Value {
		if _, ok := seenRegions[location.Name]; !ok {
			data = append(data, []string{location.Name, location.DisplayName, "false"})
		}
	}

	table := table.NewTable(table.WebApp)
	table.AppendBulk(data)
	table.Render()

	return nil
}

func init() {
	rootCmd.AddCommand(webAppCmd)

	webAppCmd.Flags().StringP("subscription-id", "s", "", "The Azure subscription id")
	// Required
	if err := webAppCmd.MarkFlagRequired("subscription-id"); err != nil {
		webAppCmd.Printf("Error marking flag required: %s", err)
		os.Exit(1)
	}

	webAppCmd.Flags().StringArrayP("location", "l", []string{}, "The Azure location to list the capabilities. Can be specified multiple times")
	webAppCmd.Flags().Bool("all-locations", false, "Whether to list capabilities for all locations")
	webAppCmd.MarkFlagsOneRequired("location", "all-locations")
	webAppCmd.MarkFlagsMutuallyExclusive("location", "all-locations")

	webAppCmd.Flags().StringP(webAppOperatingSystemChoice.Name, "o", webAppOperatingSystemChoice.Default, webAppOperatingSystemChoice.Description)
	webAppCmd.Flags().StringP(publishType.Name, "p", publishType.Default, publishType.Description)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickstartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickstartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
