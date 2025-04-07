package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nickdala/azure-resource-verifier/internal/azure"
)

// getLocations defines the resource template and handler for getting Azure regions
func getLocations(cred *azidentity.DefaultAzureCredential, subscriptionId string) (tool mcp.Tool, handler server.ToolHandlerFunc) {

	return mcp.NewTool(
			"get-azure-regions",
			mcp.WithDescription(`Get all Azure regions in the Azure subscription. This tool is used to get a list of all Azure regions
    in the Azure subscription. The list of locations is used to verify if the resource can be deployed to an Azure region.
    If the command fails because of authentication errors, ask the user to login to Azure CLI with 'az login' command.`),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			azureLocationLocator := azure.NewAzureLocationLocator(cred, ctx, subscriptionId)
			azureLocations, err := azureLocationLocator.GetLocations()
			if err != nil {
				return nil, err
			}

			r, err := json.Marshal(azureLocations)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal response: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}
