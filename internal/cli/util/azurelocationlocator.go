package util

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

type AzureLocationLocator struct {
	cred           *azidentity.DefaultAzureCredential
	ctx            context.Context
	subscriptionId string
}

type AzureLocation struct {
	Name        string
	DisplayName string
}

type AzureLocationList struct {
	Value []*AzureLocation
}

func NewAzureLocationLocator(cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) *AzureLocationLocator {
	return &AzureLocationLocator{
		cred:           cred,
		ctx:            ctx,
		subscriptionId: subscriptionId,
	}
}

func (a *AzureLocationLocator) GetLocations() (*AzureLocationList, error) {
	clientFactory, err := armsubscriptions.NewClientFactory(a.cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	locations := &AzureLocationList{
		Value: []*AzureLocation{},
	}

	inclExtendedLocations := false
	pager := clientFactory.NewClient().NewListLocationsPager(a.subscriptionId, &armsubscriptions.ClientListLocationsOptions{IncludeExtendedLocations: &inclExtendedLocations})
	for pager.More() {
		page, err := pager.NextPage(a.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to advance page: %v", err)
		}

		for _, location := range page.Value {
			azureLocation := &AzureLocation{
				Name:        *location.Name,
				DisplayName: *location.DisplayName,
			}

			locations.Value = append(locations.Value, azureLocation)
		}
	}
	return locations, nil
}
