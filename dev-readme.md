# Azure Resource Verifier

## Development

1. Initialize
    * `go mod init github.com/nickdala/azure-verifier`
    * `go mod tidy`

1. Go Modules
    * go get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresqlflexibleservers
    * go get -u github.com/Azure/azure-sdk-for-go/sdk/azidentity

Major Version Upgrade

Go uses semantic import versioning to ensure a good backward compatibility for modules. For Azure Go management SDK, we usually upgrade module version according to the corresponding service's API version. Regarding it could be a complicated experience for major version upgrade, we will try our best to keep the SDK API stable and release new version in backward compatible way. However, if any unavoidable breaking changes and a new major version releases for SDK modules, you could use these commands under your module folder to upgrade:

```
 go install github.com/icholy/gomajor@latest
 gomajor get github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice@latest
 ```

 Reference: [armappservice](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v4#section-readme)

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