/*
Copyright © 2024 Nick Dalalelis
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v4"
	"github.com/nickdala/azure-resource-verifier/internal/cli"
	"github.com/nickdala/azure-resource-verifier/internal/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	webAppOperatingSystemChoice = cli.CliChoice{
		Name:        "operating-system",
		Description: "The operating system of the web app",
		Default:     "linux",
		Choices:     []string{"linux", "windows"},
	}

	publishType = cli.CliChoice{
		Name:        "publish-type",
		Description: "The publish type of the web app",
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

	osType := viper.GetString(webAppOperatingSystemChoice.Name)
	if valid := webAppOperatingSystemChoice.IsValidChoice(osType); !valid {
		return cli.CreateAzrErr(fmt.Sprintf("Invalid operating system choice: %s", osType), nil)
	}

	publish := viper.GetString(publishType.Name)
	if valid := publishType.IsValidChoice(publish); !valid {
		return cli.CreateAzrErr(fmt.Sprintf("Invalid publish type choice: %s", publish), nil)
	}

	azureLocations, err := getLocations(cmd, cred, ctx, subscriptionId)
	if err != nil {
		return cli.CreateAzrErr("Error parsing location flag", err)
	}

	displayNameToLocation := make(map[string]string)
	for _, location := range azureLocations {
		displayNameToLocation[location.DisplayName] = location.Name
	}

	geoRegionOptions := armappservice.WebSiteManagementClientListGeoRegionsOptions{}
	if osType == "linux" {
		geoRegionOptions.LinuxWorkersEnabled = to.Ptr(true)
	}

	if publish == "container" && osType == "windows" {
		geoRegionOptions.XenonWorkersEnabled = to.Ptr(true)
	}

	seenRegions := make(map[string]struct{})

	var data [][]string

	clientFactory, err := armappservice.NewClientFactory(subscriptionId, cred, nil)
	if err != nil {
		return cli.CreateAzrErr("failed to create the app service client factory", err)
	}

	webSiteManagementClient := clientFactory.NewWebSiteManagementClient()
	pager := webSiteManagementClient.NewListGeoRegionsPager(&geoRegionOptions)

	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			break
		}

		for _, geoRegion := range nextResult.Value {
			if regionName, ok := displayNameToLocation[*geoRegion.Properties.DisplayName]; ok {
				data = append(data, []string{regionName, *geoRegion.Properties.DisplayName, "true"})
				seenRegions[*geoRegion.Properties.DisplayName] = struct{}{}
			}
		}
	}

	// Now add the regions that were not returned by the API
	for regionDisplayName, regionName := range displayNameToLocation {
		if _, ok := seenRegions[regionDisplayName]; !ok {
			data = append(data, []string{regionName, regionDisplayName, "false"})
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
	// Required
	if err := webAppCmd.MarkFlagRequired(webAppOperatingSystemChoice.Name); err != nil {
		webAppCmd.Printf("Error marking flag required: %s", err)
		os.Exit(1)
	}

	webAppCmd.Flags().StringP(publishType.Name, "p", publishType.Default, publishType.Description)
	// Required
	if err := webAppCmd.MarkFlagRequired(publishType.Name); err != nil {
		webAppCmd.Printf("Error marking flag required: %s", err)
		os.Exit(1)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// quickstartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// quickstartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
