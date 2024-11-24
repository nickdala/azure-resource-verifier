package cmd

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azrazure "github.com/nickdala/azure-resource-verifier/internal/azure"
	"github.com/spf13/cobra"
)

// This function is used to get the locations from the command line flags or from the Azure subscription
// if the --location flag is not provided. If the --location flag is provided, the locations are filtered
// based on the locations provided in the flag.
func getLocations(cmd *cobra.Command, cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) ([]*azrazure.AzureLocation, error) {

	azureLocations, err := getAllLocationsFromSubscription(cred, ctx, subscriptionId)
	if err != nil {
		return nil, err
	}

	locations, err := cmd.Flags().GetStringArray("location")
	if err != nil {
		return nil, err
	}

	// If the --location flag is not provided, return all locations
	if len(locations) == 0 {
		return azureLocations, nil
	}

	// Filter the locations based on the locations provided in the --location flag
	// 1. Create a set of locations
	locationSet := make(map[string]struct{})
	for _, location := range locations {
		locationSet[location] = struct{}{}
	}

	// 2. Create the filtered locations list
	filteredLocations := make([]*azrazure.AzureLocation, 0)

	// 3. Iterate through the locations and add the location to the location set
	for _, location := range azureLocations {
		if _, ok := locationSet[location.Name]; ok {
			filteredLocations = append(filteredLocations, location)
		}
	}

	return filteredLocations, nil
}

// This function is used to get all the locations from the Azure subscription
func getAllLocationsFromSubscription(cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) ([]*azrazure.AzureLocation, error) {
	azureLocationLocator := azrazure.NewAzureLocationLocator(cred, ctx, subscriptionId)
	azureLocations, err := azureLocationLocator.GetLocations()
	if err != nil {
		return nil, err
	}

	return azureLocations.Value, nil
}
