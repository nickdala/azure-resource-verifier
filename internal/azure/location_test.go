package azure

import (
	"testing"
)

func TestAzureLocationList_Intersection(t *testing.T) {
	type fields struct {
		Value []*AzureLocation
	}
	type args struct {
		other *AzureLocationList
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*AzureLocation
	}{
		{
			name: "Test AzureLocationList_Intersection",
			fields: fields{
				Value: []*AzureLocation{
					{
						Name:        "norway",
						DisplayName: "Norway",
					},
					{
						Name:        "eastus",
						DisplayName: "East US",
					},
					{
						Name:        "westus",
						DisplayName: "West US",
					},
				},
			},
			args: args{
				other: &AzureLocationList{
					Value: []*AzureLocation{
						{
							Name:        "eastus",
							DisplayName: "East US",
						},
						{
							Name:        "westus",
							DisplayName: "West US",
						},
						{
							Name:        "centralus",
							DisplayName: "Central US",
						},
					},
				},
			},
			want: []*AzureLocation{
				{
					Name:        "eastus",
					DisplayName: "East US",
				},
				{
					Name:        "westus",
					DisplayName: "West US",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &AzureLocationList{
				Value: tt.fields.Value,
			}
			if got := list.Intersection(tt.args.other); len(got.Value) != len(tt.want) {
				t.Errorf("AzureLocationList.Intersection() = %v, want %v", got, tt.want)
			}
		})
	}
}
