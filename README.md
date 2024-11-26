# Azure Resource Verifier

This tool is used to verify Azure resources can be deployed to a region in a subscription. It is intended to be used as a pre-deployment check.

## Features

The following features are supported:

- Get a list of regions that are available in a subscription
- Verify that Azure Cache for Redis can be deployed to a region
- Verify that Azure Database for PostgreSQL Flexible Server can be deployed to a region
- Verify that Azure App Service can be deployed to a region

## Prerequisites

- [Go 1.23](https://golang.org/dl/)
- [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli)


## Setup

1. Install the Azure CLI
2. Run `az login` to authenticate with Azure
3. Run `az account set --subscription <subscription-id>` to set the subscription to use

## Build

```
go build
```

## Usage

```
./azure-resource-verifier quickstart -s <subscription-id> -l eastus2 -l westus3
```