package azure

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v4"
)

type AzureAppService struct {
	cred           *azidentity.DefaultAzureCredential
	ctx            context.Context
	subscriptionId string
}

type AppServiceOS int

const (
	Windows AppServiceOS = iota
	Linux
)

// Implement string method for AppServiceOS
func (os AppServiceOS) String() string {
	return [...]string{"Windows", "Linux"}[os]
}

// Implement AppServiceOSFromString method
func AppServiceOSFromString(os string) (AppServiceOS, error) {
	switch os {
	case "windows":
		return Windows, nil
	case "linux":
		return Linux, nil
	default:
		return -1, fmt.Errorf("invalid AppServiceOS: %s", os)
	}
}

type AppServicePublishType int

const (
	Code AppServicePublishType = iota
	Container
)

// Implement string method for AppServicePublishType
func (pt AppServicePublishType) String() string {
	return [...]string{"Code", "Container"}[pt]
}

// Implement AppServicePublishTypeFromString method
func AppServicePublishTypeFromString(publishType string) (AppServicePublishType, error) {
	switch publishType {
	case "code":
		return Code, nil
	case "container":
		return Container, nil
	default:
		return -1, fmt.Errorf("invalid AppServicePublishType: %s", publishType)
	}
}

func NewAzureAppService(cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) *AzureAppService {
	return &AzureAppService{
		cred:           cred,
		ctx:            ctx,
		subscriptionId: subscriptionId,
	}
}

func (a *AzureAppService) GetAppServiceLocations(locations *AzureLocationList, os AppServiceOS, publishType AppServicePublishType) (*AzureLocationList, error) {
	geoRegionOptions := armappservice.WebSiteManagementClientListGeoRegionsOptions{}

	if os == Linux {
		geoRegionOptions.LinuxWorkersEnabled = to.Ptr(true)
	}

	if publishType == Container && os == Windows {
		geoRegionOptions.XenonWorkersEnabled = to.Ptr(true)
	}

	clientFactory, err := armappservice.NewClientFactory(a.subscriptionId, a.cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create the app service client factory %w", err)
	}

	webSiteManagementClient := clientFactory.NewWebSiteManagementClient()
	pager := webSiteManagementClient.NewListGeoRegionsPager(&geoRegionOptions)

	// Create a map of display name to location
	// This is because the API returns the location display name
	displayNameToLocation := make(map[string]*AzureLocation)
	for _, location := range locations.Value {
		displayNameToLocation[location.DisplayName] = location
	}

	appServicelocations := &AzureLocationList{
		Value: []*AzureLocation{},
	}

	for pager.More() {
		nextResult, err := pager.NextPage(a.ctx)
		if err != nil {
			log.Printf("failed to get the app service locations %v", err)
		}

		for _, geoRegion := range nextResult.Value {
			if region, ok := displayNameToLocation[*geoRegion.Properties.DisplayName]; ok {
				appServicelocations.Value = append(appServicelocations.Value, region)
			} else {
				log.Printf("Location %s not found in the location list", *geoRegion.Properties.DisplayName)
			}
		}
	}

	return appServicelocations, nil
}
