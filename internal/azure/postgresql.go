package azure

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers"
)

type AzurePostgresqlFlexibleServer struct {
	cred           *azidentity.DefaultAzureCredential
	ctx            context.Context
	subscriptionId string
}

func NewAzurePostgresqlFlexibleServer(cred *azidentity.DefaultAzureCredential, ctx context.Context, subscriptionId string) *AzurePostgresqlFlexibleServer {
	return &AzurePostgresqlFlexibleServer{
		cred:           cred,
		ctx:            ctx,
		subscriptionId: subscriptionId,
	}
}

type AzurePostgresqlLocationList AzureLocationList
type AzurePostgresqlHaLocationList AzureLocationList

type AzurePostgresqlNonDeployableLocation struct {
	Name        string
	DisplayName string
	Reason      string
}
type AzurePostgresqlNonDeployableLocationList struct {
	Value []*AzurePostgresqlNonDeployableLocation
}

func (a *AzurePostgresqlFlexibleServer) GetPostgresqlLocations(locations *AzureLocationList) (*AzurePostgresqlLocationList, *AzurePostgresqlHaLocationList, *AzurePostgresqlNonDeployableLocationList, error) {
	client, err := armpostgresqlflexibleservers.NewLocationBasedCapabilitiesClient(a.subscriptionId, a.cred, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create the postgresql flexible server client %w", err)
	}

	// The following is used to store locations from our go routine.
	// We will merge the results after all go routines are done.
	deployableLocations := make([]*AzureLocation, len(locations.Value))
	haLocations := make([]*AzureLocation, len(locations.Value))
	nonDeployableLocations := make([]*AzurePostgresqlNonDeployableLocation, len(locations.Value))

	var wg sync.WaitGroup
	for i, location := range locations.Value {
		wg.Add(1)
		go func(idx int, azureLocation *AzureLocation) {
			defer wg.Done()
			pager := client.NewExecutePager(azureLocation.Name, nil)
			log.Printf("Getting capabilities for location %s", azureLocation.DisplayName)
			for pager.More() {
				nextResult, err := pager.NextPage(a.ctx)
				if err != nil {
					if azureErr, ok := err.(*azcore.ResponseError); ok {
						nonDeployableLocations[idx] =
							&AzurePostgresqlNonDeployableLocation{
								Name:        location.Name,
								DisplayName: location.DisplayName,
								Reason:      azureErr.ErrorCode,
							}
					} else {
						nonDeployableLocations[idx] =
							&AzurePostgresqlNonDeployableLocation{
								Name:        location.Name,
								DisplayName: location.DisplayName,
								Reason:      err.Error(),
							}
					}
					break
				}

				if len(nextResult.Value) == 0 {
					nonDeployableLocations[idx] =
						&AzurePostgresqlNonDeployableLocation{
							Name:        location.Name,
							DisplayName: location.DisplayName,
							Reason:      "can't deploy to this location",
						}
					break
				}

				// We have the capabilities for the location.
				// You can at least deploy PostgreSQL Flexible Server to this location.
				// Check if the location supports HA.
				haEnabled := false
				for _, capability := range nextResult.Value {
					if *capability.ZoneRedundantHaSupported {
						haLocations[idx] = location
						haEnabled = true
						break // Only need confirmation for one capability for HA
					}
				}

				// HA is not supported in this location. Add to the deployable list.
				if !haEnabled {
					deployableLocations[idx] = location
					break
				}
			}
		}(i, location)
	}

	wg.Wait()

	// Remove nil values from the deployable and ha locations.
	deployableLocations = removeNilItems(deployableLocations)
	haLocations = removeNilItems(haLocations)
	nonDeployableLocations = removeNilItems(nonDeployableLocations)

	deployableLocationList := &AzurePostgresqlLocationList{
		Value: deployableLocations,
	}

	haLocationList := &AzurePostgresqlHaLocationList{
		Value: haLocations,
	}

	nonDeployableLocationList := &AzurePostgresqlNonDeployableLocationList{
		Value: nonDeployableLocations,
	}

	return deployableLocationList, haLocationList, nonDeployableLocationList, nil
}

func removeNilItems[T any](items []*T) []*T {
	var result []*T
	for _, item := range items {
		if item != nil {
			result = append(result, item)
		}
	}
	return result
}
