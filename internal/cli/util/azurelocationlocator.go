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

func NewAzureLocationLocator(cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) *AzureLocationLocator {
	return &AzureLocationLocator{
		cred:           cred,
		ctx:            ctx,
		subscriptionId: subscriptionId,
	}
}

func (a *AzureLocationLocator) GetLocations() ([]string, error) {
	clientFactory, err := armsubscriptions.NewClientFactory(a.cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	locations := []string{}

	pager := clientFactory.NewClient().NewListLocationsPager(a.subscriptionId, &armsubscriptions.ClientListLocationsOptions{IncludeExtendedLocations: nil})
	for pager.More() {
		page, err := pager.NextPage(a.ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to advance page: %v", err)
		}

		for _, location := range page.Value {
			locations = append(locations, *location.Name)
		}
	}
	return locations, nil
}
