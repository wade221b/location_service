package constants

const (
	DynamicConsumerAPI = "v1/venues/%s/dynamic"
	StaticConsumerAPI  = "v1/venues/%s/static"

	//ENV vars
	PricingServiceHost = "PRICING_SERVICE_HOST"
	ConsumerApi        = "CONSUMER_API"
)

const (
	LOCATION_SERVICE_WRONG_LOCATION_TYPE_400   = "LOCATION_SERVICE_WRONG_LOCATION_TYPE"
	LOCATION_SERVICE_DISTANCE_OUT_OF_REACH_400 = "LOCATION_SERVICE_DISTANCE_OUT_OF_REACH_400"
	LOCATION_SERVICE_INTERNAL_SERVER_ERROR_500 = "LOCATION_SERVICE_INTERNAL_SERVER_ERROR_500"
	// Add more error codes as needed...
)

// Map error codes to their default messages.
var ErrorMessages = map[string]string{
	LOCATION_SERVICE_WRONG_LOCATION_TYPE_400:   "The provided location type is incorrect.",
	LOCATION_SERVICE_DISTANCE_OUT_OF_REACH_400: "The customer distance is out of reach",
	LOCATION_SERVICE_INTERNAL_SERVER_ERROR_500: "Something went wrong. ",
}
