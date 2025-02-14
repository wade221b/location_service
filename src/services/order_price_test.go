package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/your-username/your-project/src/constants"
	mock_httpclient "github.com/your-username/your-project/src/mocks"
	reqResp "github.com/your-username/your-project/src/req_resp"
	// "reflect"
	// "testing"
	// reqResp "github.com/your-username/your-project/src/req_resp"
	// "fmt"
	// "reflect"
	// "testing"
	// "myproject/constants"
	// "myproject/delivery"
	// "myproject/mock_httpclient" // Generated mock package.
	// "myproject/reqResp"
	// "myproject/utils"
	// "github.com/golang/mock/gomock"
)

// func Test_deliveryPriceCalculator_CalculateDeliveryOrderPrice(t *testing.T) {
// 	type fields struct {
// 		client *httpclient.HTTPClient
// 	}
// 	type args struct {
// 		venueSlug string
// 		cartValue int
// 		userLat   float64
// 		userLon   float64
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    *reqResp.DeliveryResponse
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			d := &deliveryPriceCalculator{
// 				client: tt.fields.client,
// 			}
// 			got, err := d.CalculateDeliveryOrderPrice(tt.args.venueSlug, tt.args.cartValue, tt.args.userLat, tt.args.userLon)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("deliveryPriceCalculator.CalculateDeliveryOrderPrice() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("deliveryPriceCalculator.CalculateDeliveryOrderPrice() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// import (
// 	"encoding/json"
// 	"fmt"
// 	"math"
// 	"reflect"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"myproject/constants"
// 	"myproject/delivery"
// 	"myproject/mock_httpclient" // Generated mock package.
// 	"myproject/reqResp"
// 	"myproject/utils"
// )

// For testing, we define the handle* functions as used by the calculator.
// In your actual code these might be defined elsewhere.
// func handleStaticResponseData(b []byte) (*reqResp.StaticResponseConsumerAPI, error) {
// 	var resp reqResp.StaticResponseConsumerAPI
// 	err := json.Unmarshal(b, &resp)
// 	return &resp, err
// }

// func handleDynamicResponseData(b []byte) (*reqResp.DynamicResponseConsumerAPI, error) {
// 	var resp reqResp.DynamicResponseConsumerAPI
// 	err := json.Unmarshal(b, &resp)
// 	return &resp, err
// }

func Test_deliveryPriceCalculator_CalculateDeliveryOrderPrice(t *testing.T) {
	type args struct {
		venueSlug string
		cartValue int
		userLat   float64
		userLon   float64
	}
	tests := []struct {
		name      string
		setupMock func(m *mock_httpclient.MockClient, t *testing.T)
		args      args
		want      *reqResp.DeliveryResponse
		wantErr   bool
	}{
		{
			name: "successful calculation with no surcharge",
			setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
				// Expect call for static API
				m.EXPECT().
					Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue1")).
					Return([]byte(`{
						"VenueRaw": {
							"Location": {
								"Coordinates": [0.0, 0.0]
							}
						}
					}`), nil).Times(1)
				// Expect call for dynamic API
				m.EXPECT().
					Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue1")).
					Return([]byte(`{
						"VenueRaw": {
							"DeliverySpecs": {
								"OrderMinimumNoSurcharge": 100,
								"DeliveryPricing": {
									"BasePrice": 50
								}
							}
						}
					}`), nil).Times(1)
			},
			args: args{
				venueSlug: "venue1",
				cartValue: 100,
				userLat:   0.0,
				userLon:   0.0,
			},
			// Calculation:
			// - Distance = 0
			// - Fee = BasePrice (50) + A (10) + B*0/10 = 60
			// - TotalPrice = cartValue + fee = 100 + 60 = 160
			// - Surcharge = 0 since cartValue >= OrderMinimumNoSurcharge
			want: &reqResp.DeliveryResponse{
				TotalPrice:          160,
				SmallOrderSurcharge: 0,
				CartValue:           100,
				Delivery: reqResp.DeliveryDetail{
					Fee:      60,
					Distance: 0,
				},
			},
			wantErr: false,
		},
		// {
		// 	name: "cart value below order minimum (with surcharge)",
		// 	setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
		// 		m.EXPECT().
		// 			Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue2")).
		// 			Return([]byte(`{
		// 				"VenueRaw": {
		// 					"Location": {
		// 						"Coordinates": [0.0, 0.0]
		// 					}
		// 				}
		// 			}`), nil).Times(1)
		// 		m.EXPECT().
		// 			Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue2")).
		// 			Return([]byte(`{
		// 				"VenueRaw": {
		// 					"DeliverySpecs": {
		// 						"OrderMinimumNoSurcharge": 100,
		// 						"DeliveryPricing": {
		// 							"BasePrice": 50
		// 						}
		// 					}
		// 				}
		// 			}`), nil).Times(1)
		// 	},
		// 	args: args{
		// 		venueSlug: "venue2",
		// 		cartValue: 80,
		// 		userLat:   0.0,
		// 		userLon:   0.0,
		// 	},
		// 	// Calculation:
		// 	// - Fee remains 60.
		// 	// - TotalPrice = 80 + 60 = 140.
		// 	// - Surcharge = 100 - 80 = 20.
		// 	want: &reqResp.DeliveryResponse{
		// 		TotalPrice:          140,
		// 		SmallOrderSurcharge: 20,
		// 		CartValue:           80,
		// 		Delivery: reqResp.DeliveryDetail{
		// 			Fee:      60,
		// 			Distance: 0,
		// 		},
		// 	},
		// 	wantErr: false,
		// },
		// {
		// 	name: "distance out of reach returns error",
		// 	setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
		// 		// Override utils.GetDistanceRange to simulate an out-of-reach scenario.
		// 		// originalGetDistanceRange := utils.GetDistanceRange
		// 		// utils.GetDistanceRange = func(dynamic *reqResp.DynamicResponseConsumerAPI, distance float64) *utils.GetDistanceRange {
		// 		// 	return nil
		// 		// }
		// 		// t.Cleanup(func() {
		// 		// 	utils.GetDistanceRange = originalGetDistanceRange
		// 		// })
		// 		m.EXPECT().
		// 			Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue3")).
		// 			Return([]byte(`{
		// 				"VenueRaw": {
		// 					"Location": {
		// 						"Coordinates": [10.0, 10.0]
		// 					}
		// 				}
		// 			}`), nil).Times(1)
		// 		m.EXPECT().
		// 			Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue3")).
		// 			Return([]byte(`{
		// 				"VenueRaw": {
		// 					"DeliverySpecs": {
		// 						"OrderMinimumNoSurcharge": 100,
		// 						"DeliveryPricing": {
		// 							"BasePrice": 50
		// 						}
		// 					}
		// 				}
		// 			}`), nil).Times(1)
		// 	},
		// 	args: args{
		// 		venueSlug: "venue3",
		// 		cartValue: 100,
		// 		userLat:   0.0,
		// 		userLon:   0.0,
		// 	},
		// 	want:    nil,
		// 	wantErr: true,
		// },
		// {
		// 	name: "error in static API response",
		// 	setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
		// 		m.EXPECT().
		// 			Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue4")).
		// 			Return(nil, fmt.Errorf("static API error")).Times(1)
		// 	},
		// 	args: args{
		// 		venueSlug: "venue4",
		// 		cartValue: 100,
		// 		userLat:   0.0,
		// 		userLon:   0.0,
		// 	},
		// 	want:    nil,
		// 	wantErr: true,
		// },
		// {
		// 	name: "error in dynamic API response",
		// 	setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
		// 		m.EXPECT().
		// 			Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue5")).
		// 			Return([]byte(`{
		// 				"VenueRaw": {
		// 					"Location": {
		// 						"Coordinates": [0.0, 0.0]
		// 					}
		// 				}
		// 			}`), nil).Times(1)
		// 		m.EXPECT().
		// 			Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue5")).
		// 			Return(nil, fmt.Errorf("dynamic API error")).Times(1)
		// 	},
		// 	args: args{
		// 		venueSlug: "venue5",
		// 		cartValue: 100,
		// 		userLat:   0.0,
		// 		userLon:   0.0,
		// 	},
		// 	want:    nil,
		// 	wantErr: true,
		// },
	}

	// Run each test case.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mock_httpclient.NewMockClient(ctrl)
			tt.setupMock(mockClient, t)

			// Create the delivery calculator using the mock client.

			calculator := NewDeliveryPriceCalculator(mockClient)
			fmt.Println("got the calculator")
			got, err := calculator.CalculateDeliveryOrderPrice(tt.args.venueSlug, tt.args.cartValue, tt.args.userLat, tt.args.userLon)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateDeliveryOrderPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateDeliveryOrderPrice() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
