package cmd

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/nickdala/azure-resource-verifier/internal/cli/util"
	"github.com/spf13/cobra"
)

func getLocations(cmd *cobra.Command, cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) ([]string, error) {
	locations, err := cmd.Flags().GetStringArray("location")
	if err != nil {
		return nil, err
	}

	if len(locations) > 0 {
		return locations, nil
	}

	azureLocationLocator := util.NewAzureLocationLocator(cred, ctx, subscriptionId)
	azureLocations, err := azureLocationLocator.GetLocations()
	if err != nil {
		return nil, err
	}

	azureLocationsLength := len(azureLocations.Value)
	locations = make([]string, azureLocationsLength)
	for i, location := range azureLocations.Value {
		locations[i] = location.Name
	}

	return locations, nil
}
