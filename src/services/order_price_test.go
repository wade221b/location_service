package service

import (
	"reflect"
	"testing"

	httpclient "github.com/your-username/your-project/src/client"
	reqResp "github.com/your-username/your-project/src/req_resp"
)

func Test_deliveryPriceCalculator_CalculateDeliveryOrderPrice(t *testing.T) {
	type fields struct {
		client *httpclient.HTTPClient
	}
	type args struct {
		venueSlug string
		cartValue int
		userLat   float64
		userLon   float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *reqResp.DeliveryResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &deliveryPriceCalculator{
				client: tt.fields.client,
			}
			got, err := d.CalculateDeliveryOrderPrice(tt.args.venueSlug, tt.args.cartValue, tt.args.userLat, tt.args.userLon)
			if (err != nil) != tt.wantErr {
				t.Errorf("deliveryPriceCalculator.CalculateDeliveryOrderPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deliveryPriceCalculator.CalculateDeliveryOrderPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}
