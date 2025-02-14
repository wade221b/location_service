package constants

const (
	DynamicConsumerAPI = "v1/venues/%s/dynamic"
	StaticConsumerAPI  = "v1/venues/%s/static"

	//ENV vars
	PricingServiceHost = "PRICING_SERVICE_HOST"
	ConsumerApi        = "CONSUMER_API"
)

const (
    CS_SERVICE_WRONG_LOCATION_TYPE = "CS_SERVICE_WRONG_LOCATION_TYPE"
    // Add more error codes as needed...
)

// Map error codes to their default messages.
var ErrorMessages = map[string]string{
    CS_SERVICE_WRONG_LOCATION_TYPE: "The provided location type is incorrect.",
    
}