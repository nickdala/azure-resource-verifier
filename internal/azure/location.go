package azure

type AzureLocation struct {
	Name        string
	DisplayName string
}

type AzureLocationList struct {
	Value []*AzureLocation
}

func (list *AzureLocationList) Intersection(other *AzureLocationList) *AzureLocationList {
	seenRegions := make(map[string]struct{})
	var data []*AzureLocation

	for _, location := range list.Value {
		seenRegions[location.Name] = struct{}{}
	}

	for _, location := range other.Value {
		if _, ok := seenRegions[location.Name]; ok {
			data = append(data, location)
		}
	}

	return &AzureLocationList{Value: data}
}

func (list *AzureLocationList) Difference(other *AzureLocationList) *AzureLocationList {
	var data []*AzureLocation

	otherRegions := make(map[string]struct{})
	for _, location := range other.Value {
		otherRegions[location.Name] = struct{}{}
	}

	for _, location := range list.Value {
		if _, ok := otherRegions[location.Name]; !ok {
			data = append(data, location)
		}
	}

	return &AzureLocationList{Value: data}
}
