package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
)

type AzureRedisCache struct {
	cred           *azidentity.DefaultAzureCredential
	ctx            context.Context
	subscriptionId string
}

func NewAzureRedisCache(cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) *AzureRedisCache {
	return &AzureRedisCache{
		cred:           cred,
		ctx:            ctx,
		subscriptionId: subscriptionId,
	}
}

func (a *AzureRedisCache) GetRedisLocations() (*AzureLocationList, error) {
	clientFactory, err := armresources.NewClientFactory(a.subscriptionId, a.cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the arm resource client factory %w", err)
	}

	res, err := clientFactory.NewProvidersClient().Get(a.ctx, "Microsoft.Cache", &armresources.ProvidersClientGetOptions{Expand: nil})
	if err != nil {
		return nil, fmt.Errorf("failed to get the cache provider %w", err)
	}

	azureLocationLocator := NewAzureLocationLocator(a.cred, a.ctx, a.subscriptionId)
	azureLocations, err := azureLocationLocator.GetLocations()
	if err != nil {
		return nil, err
	}

	// Create a map of display name to location
	displayNameToLocation := make(map[string]*AzureLocation)
	for _, location := range azureLocations.Value {
		displayNameToLocation[location.DisplayName] = location
	}

	redisLocations, err := getRedisLocations(&res.Provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get the Azure Cache for Redis locations %w", err)
	}

	locations := &AzureLocationList{
		Value: []*AzureLocation{},
	}

	for _, location := range redisLocations {
		if regionName, ok := displayNameToLocation[*location]; ok {
			locations.Value = append(locations.Value, regionName)
		}
	}

	return locations, nil
}

func getRedisLocations(provider *armresources.Provider) ([]*string, error) {
	if provider.ResourceTypes == nil {
		return nil, fmt.Errorf("failed to get the cache provider resource types")
	}

	for _, resourceType := range provider.ResourceTypes {
		if resourceType.ResourceType == nil {
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

	return nil, fmt.Errorf("no Redis locations found")
}
