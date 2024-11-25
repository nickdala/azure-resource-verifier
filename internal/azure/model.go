package azure

type AzureLocation struct {
	Name        string
	DisplayName string
}

type AzureLocationList struct {
	Value []*AzureLocation
}
