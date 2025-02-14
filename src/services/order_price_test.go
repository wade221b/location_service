package service

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/your-username/your-project/src/constants"
	mock_httpclient "github.com/your-username/your-project/src/mocks"
	reqResp "github.com/your-username/your-project/src/req_resp"
)

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
			name: "successful calculation with no surcharge using single range",
			setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
				// Static API returns coordinates [0,0].
				m.EXPECT().
					Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue1")).
					Return([]byte(`{
						"venue_raw": {
							"location": {
								"coordinates": [0.0, 0.0]
							}
						}
					}`), nil).Times(1)
				// Dynamic API returns order minimum = 100, base price = 50,
				// and a single distance range where min is 0 and max is 0.
				m.EXPECT().
					Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue1")).
					Return([]byte(`{
						"venue_raw": {
							"delivery_specs": {
								"order_minimum_no_surcharge": 100,
								"delivery_pricing": {
									"base_price": 50,
									"distance_ranges": [
										{ "min": 0, "max": 0, "a": 10, "b": 2, "flag": null }
									]
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
			// For distance 0: fee = basePrice (50) + A (10) + B*(0/10) = 60.
			// TotalPrice = 100 + 60 = 160, surcharge = 0.
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
		{
			name: "cart value below order minimum (with surcharge) using single range",
			setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
				m.EXPECT().
					Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue2")).
					Return([]byte(`{
						"venue_raw": {
							"location": {
								"coordinates": [0.0, 0.0]
							}
						}
					}`), nil).Times(1)
				m.EXPECT().
					Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue2")).
					Return([]byte(`{
						"venue_raw": {
							"delivery_specs": {
								"order_minimum_no_surcharge": 100,
								"delivery_pricing": {
									"base_price": 50,
									"distance_ranges": [
										{ "min": 0, "max": 0, "a": 10, "b": 2, "flag": null }
									]
								}
							}
						}
					}`), nil).Times(1)
			},
			args: args{
				venueSlug: "venue2",
				cartValue: 80, // below the minimum of 100
				userLat:   0.0,
				userLon:   0.0,
			},
			// Fee is still 60. Surcharge = 100 - 80 = 20, TotalPrice = 80 + 60 = 140.
			want: &reqResp.DeliveryResponse{
				TotalPrice:          140,
				SmallOrderSurcharge: 20,
				CartValue:           80,
				Delivery: reqResp.DeliveryDetail{
					Fee:      60,
					Distance: 0,
				},
			},
			wantErr: false,
		},
		{
			name: "multiple distance ranges with second range selected",
			setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
				// Static API returns coordinates [10,10]. With user at [0,0], distance ≈ 14.14.
				m.EXPECT().
					Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue_multi")).
					Return([]byte(`{
						"venue_raw": {
							"location": {
								"coordinates": [10.0, 10.0]
							}
						}
					}`), nil).Times(1)
				// Dynamic API returns two distance ranges:
				//  - First range: from 0 to 5 (min=0, max=5) with a=10, b=2.
				//  - Second range: from 5 to 0 (max=0 indicates no upper limit) with a=15, b=3.
				// For a distance ≈ 14.14, the second range should be selected.
				m.EXPECT().
					Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue_multi")).
					Return([]byte(`{
						"venue_raw": {
							"delivery_specs": {
								"order_minimum_no_surcharge": 100,
								"delivery_pricing": {
									"base_price": 50,
									"distance_ranges": [
										{ "min": 0, "max": 5, "a": 10, "b": 2, "flag": null },
										{ "min": 5, "max": 0, "a": 15, "b": 3, "flag": null }
									]
								}
							}
						}
					}`), nil).Times(1)
			},
			args: args{
				venueSlug: "venue_multi",
				cartValue: 100,
				userLat:   0.0,
				userLon:   0.0,
			},
			// For distance ~14.14, second range is used:
			// fee = basePrice (50) + second range's a (15) + b*(distance/10) = 50 + 15 + 3*(14.14/10)
			//    = 50 + 15 + 3*1.414 ≈ 50 + 15 + 4.242 = 69.242 => rounds to 69.
			// TotalPrice = cartValue (100) + fee (69) = 169.
			// Delivery distance is math.Round(14.14) = 14.
			want: &reqResp.DeliveryResponse{
				TotalPrice:          169,
				SmallOrderSurcharge: 0,
				CartValue:           100,
				Delivery: reqResp.DeliveryDetail{
					Fee:      69,
					Distance: 14,
				},
			},
			wantErr: false,
		},
		{
			name: "distance out of reach returns error",
			setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
				// Override utils.GetDistanceRange to simulate an out-of-reach scenario.
				// originalGetDistanceRange := utils.GetDistanceRange
				// utils.GetDistanceRange = func(dynamic *reqResp.DynamicResponseConsumerAPI, distance float64) *utils.DistanceRange {
				// 	return nil
				// }
				// t.Cleanup(func() {
				// 	utils.GetDistanceRange = originalGetDistanceRange
				// })
				// Provide static coordinates that yield a non-zero distance.
				m.EXPECT().
					Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue3")).
					Return([]byte(`{
						"venue_raw": {
							"location": {
								"coordinates": [10.0, 10.0]
							}
						}
					}`), nil).Times(1)
				m.EXPECT().
					Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue3")).
					Return([]byte(`{
						"venue_raw": {
							"delivery_specs": {
								"order_minimum_no_surcharge": 100,
								"delivery_pricing": {
									"base_price": 50,
									"distance_ranges": [
										{ "min": 0, "max": 0, "a": 10, "b": 2, "flag": null }
									]
								}
							}
						}
					}`), nil).Times(1)
			},
			args: args{
				venueSlug: "venue3",
				cartValue: 100,
				userLat:   0.0,
				userLon:   0.0,
			},
			// Since GetDistanceRange returns nil, an error is expected.
			want:    nil,
			wantErr: true,
		},
		{
			name: "error in static API response",
			setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
				m.EXPECT().
					Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue4")).
					Return(nil, fmt.Errorf("static API error")).Times(1)
			},
			args: args{
				venueSlug: "venue4",
				cartValue: 100,
				userLat:   0.0,
				userLon:   0.0,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error in dynamic API response",
			setupMock: func(m *mock_httpclient.MockClient, t *testing.T) {
				m.EXPECT().
					Get(fmt.Sprintf(constants.StaticConsumerAPI, "venue5")).
					Return([]byte(`{
						"venue_raw": {
							"location": {
								"coordinates": [0.0, 0.0]
							}
						}
					}`), nil).Times(1)
				m.EXPECT().
					Get(fmt.Sprintf(constants.DynamicConsumerAPI, "venue5")).
					Return(nil, fmt.Errorf("dynamic API error")).Times(1)
			},
			args: args{
				venueSlug: "venue5",
				cartValue: 100,
				userLat:   0.0,
				userLon:   0.0,
			},
			want:    nil,
			wantErr: true,
		},
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
