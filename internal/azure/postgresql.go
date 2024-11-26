package azure

import (
	"context"
	"fmt"
	"log"

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

	deployableLocations := &AzurePostgresqlLocationList{
		Value: []*AzureLocation{},
	}

	haLocations := &AzurePostgresqlHaLocationList{
		Value: []*AzureLocation{},
	}

	nonDeployableLocations := &AzurePostgresqlNonDeployableLocationList{
		Value: []*AzurePostgresqlNonDeployableLocation{},
	}

	for _, location := range locations.Value {
		pager := client.NewExecutePager(location.Name, nil)
		log.Printf("Getting capabilities for location %s", location.DisplayName)
		for pager.More() {
			nextResult, err := pager.NextPage(a.ctx)
			if err != nil {
				if azureErr, ok := err.(*azcore.ResponseError); ok {
					nonDeployableLocations.Value = append(nonDeployableLocations.Value,
						&AzurePostgresqlNonDeployableLocation{
							Name:        location.Name,
							DisplayName: location.DisplayName,
							Reason:      azureErr.ErrorCode,
						})
				} else {
					nonDeployableLocations.Value = append(nonDeployableLocations.Value,
						&AzurePostgresqlNonDeployableLocation{
							Name:        location.Name,
							DisplayName: location.DisplayName,
							Reason:      err.Error(),
						})
				}
				break
			}

			if len(nextResult.Value) == 0 {
				nonDeployableLocations.Value = append(nonDeployableLocations.Value,
					&AzurePostgresqlNonDeployableLocation{
						Name:        location.Name,
						DisplayName: location.DisplayName,
						Reason:      "can't deploy to this location",
					})
				break
			}

			// We have the capabilities for the location.
			// You can at least deploy PostgreSQL Flexible Server to this location.
			// Check if the location supports HA.
			haEnabled := false
			for _, capability := range nextResult.Value {
				if *capability.ZoneRedundantHaSupported {
					haLocations.Value = append(haLocations.Value, location)
					haEnabled = true
					break // Only need confirmation for one capability for HA
				}
			}

			// HA is not supported in this location. Add to the deployable list.
			if !haEnabled {
				deployableLocations.Value = append(deployableLocations.Value, location)
				break
			}
		}
	}

	return deployableLocations, haLocations, nonDeployableLocations, nil
}
