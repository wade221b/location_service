package reqResp

// DeliveryResponse represents the final JSON response structure
type DeliveryResponse struct {
	TotalPrice          int            `json:"total_price"`
	SmallOrderSurcharge int            `json:"small_order_surcharge"`
	CartValue           int            `json:"cart_value"`
	Delivery            DeliveryDetail `json:"delivery"`
}

// DeliveryDetail represents the details about the delivery cost and distance
type DeliveryDetail struct {
	Fee      int `json:"fee"`
	Distance int `json:"distance"`
}

type DynamicResponseConsumerAPI struct {
	VenueRaw VenueRaw `json:"venue_raw"`
}

type VenueRaw struct {
	DeliverySpecs DeliverySpecs `json:"delivery_specs"`
}

type DeliverySpecs struct {
	OrderMinimumNoSurcharge int             `json:"order_minimum_no_surcharge"`
	DeliveryPricing         DeliveryPricing `json:"delivery_pricing"`
}

type DeliveryPricing struct {
	BasePrice      int             `json:"base_price"`
	DistanceRanges []DistanceRange `json:"distance_ranges"`
}

type DistanceRange struct {
	Min  int         `json:"min"`
	Max  int         `json:"max"`
	A    float64     `json:"a"`
	B    float64     `json:"b"`
	Flag interface{} `json:"flag"`
	// or *string / json.RawMessage if you need more specific handling
}

type StaticResponseConsumerAPI struct {
	VenueRaw struct {
		Location struct {
			Coordinates []float64 `json:"coordinates"`
		} `json:"location"`
	} `json:"venue_raw"`
}
