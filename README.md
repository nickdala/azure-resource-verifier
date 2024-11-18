# Azure Resource Verifier

## Development

1. Initialize
    * `go mod init github.com/nickdala/azure-verifier`
    * `go mod tidy`

1. Go Modules
    * go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers
    * go get -u github.com/Azure/azure-sdk-for-go/sdk/azidentity

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