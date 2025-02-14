package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	httpclient "github.com/your-username/your-project/src/client"
	"github.com/your-username/your-project/src/constants"
	reqResp "github.com/your-username/your-project/src/req_resp"
	"github.com/your-username/your-project/src/utils"
	// reqResp "github.com/your-username/your-project/src/req_resp/req_resp"
)

type DeliveryPriceCalculator interface {
	CalculateDeliveryOrderPrice(venueSlug string, cartValue int, userLat, userLon float64) (*reqResp.DeliveryResponse, error)
}
type deliveryPriceCalculator struct {
	client httpclient.Client
}

func NewDeliveryPriceCalculator(client httpclient.Client) DeliveryPriceCalculator {
	return &deliveryPriceCalculator{
		client: client,
	}
}

func (d *deliveryPriceCalculator) getAPIResponse(hostString, venue string) ([]byte, error) {
	host := fmt.Sprintf(hostString, venue)
	bytes, err := d.client.Get(host)

	if err != nil {
		err = fmt.Errorf(err.Error())
		log.Printf("Error in getAPIResponse: %v", err)
		return nil, err
	}
	return bytes, nil

	// dd, err := handleDynamicResponseData(bytes)
}

func (d *deliveryPriceCalculator) getDynamicVenueDetails(venue string) (*reqResp.DynamicResponseConsumerAPI, error) {

	bytes, err := d.getAPIResponse(constants.DynamicConsumerAPI, venue)
	// staticHost := fmt.Sprintf(constants.DynamicConsumerAPI, venue)
	// bytes, err := d.client.Get(dynamicHost)

	if err != nil {
		err = fmt.Errorf(err.Error())
		log.Printf("Error in getDynamicVenueDetails: %v", err)
		return nil, err
	}

	dd, err := handleDynamicResponseData(bytes)
	if err != nil {
		err = fmt.Errorf(err.Error())
		log.Printf("Error in getDynamicVenueDetails: %v", err)
		return nil, err
	}
	return dd, nil
}
func (d *deliveryPriceCalculator) getStaticVenueDetails(venue string) (*reqResp.StaticResponseConsumerAPI, error) {

	bytes, err := d.getAPIResponse(constants.StaticConsumerAPI, venue)
	// staticHost := fmt.Sprintf(constants.DynamicConsumerAPI, venue)
	// bytes, err := d.client.Get(dynamicHost)

	if err != nil {
		err = fmt.Errorf(err.Error())
		log.Printf("Error in getStaticVenueDetails: %v", err)
		return nil, err
	}

	dd, err := handleStaticResponseData(bytes)
	if err != nil {
		err = fmt.Errorf(err.Error())
		log.Printf("Error in getStaticVenueDetails: %v", err)
		return nil, err
	}
	return dd, nil
}

// CalculateDeliveryOrderPrice calculates the entire cost breakdown of a delivery order
func (d *deliveryPriceCalculator) CalculateDeliveryOrderPrice(venueSlug string, cartValue int, userLat, userLon float64) (*reqResp.DeliveryResponse, error) {

	var minOrder int

	staticVenueDetails, err := d.getStaticVenueDetails(venueSlug) //get static values for venue
	if err != nil {
		err = fmt.Errorf(err.Error())
		log.Printf("Error in CalculateDeliveryOrderPrice: %v", err)
		return nil, err
	}

	fmt.Println("got static venud details staticVenueDetails ", staticVenueDetails)

	dynamicVenueDetails, err := d.getDynamicVenueDetails(venueSlug) //get dynamic values for venue
	if err != nil {
		err = fmt.Errorf(err.Error())
		log.Printf("Error in CalculateDeliveryOrderPrice: %v", err)
		return nil, err
	}

	fmt.Println("got dynamic venud details staticVenueDetails ", dynamicVenueDetails)

	log.Printf("order_minimum_no_surcharge = %v, cartvalue = %v",
		dynamicVenueDetails.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge, cartValue)

	longitude := staticVenueDetails.VenueRaw.Location.Coordinates[0]
	latitude := staticVenueDetails.VenueRaw.Location.Coordinates[1]
	distance := utils.CalculateExactDistance(latitude, longitude, userLat, userLon)
	minOrder = dynamicVenueDetails.VenueRaw.DeliverySpecs.OrderMinimumNoSurcharge
	basePrice := dynamicVenueDetails.VenueRaw.DeliverySpecs.DeliveryPricing.BasePrice

	smallOrderSurcharge := 0
	if cartValue < minOrder { //todo : check the int float thing here
		smallOrderSurcharge = minOrder - cartValue
	}

	distanceRange := utils.GetDistanceRange(dynamicVenueDetails, distance)
	if distanceRange == nil {
		return nil, fmt.Errorf("distance out of reach")
	}

	deliveryFee := float64(basePrice) + distanceRange.A + distanceRange.B*distance/10
	totalPrice := cartValue + int(math.Round(deliveryFee))

	return &reqResp.DeliveryResponse{
		TotalPrice:          totalPrice,
		SmallOrderSurcharge: smallOrderSurcharge,
		CartValue:           cartValue,
		Delivery: reqResp.DeliveryDetail{
			Fee:      int(math.Round(deliveryFee)),
			Distance: int(math.Round(distance)),
		},
	}, nil
}

func handleDynamicResponseData(body []byte) (*reqResp.DynamicResponseConsumerAPI, error) {
	var resp reqResp.DynamicResponseConsumerAPI
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
func handleStaticResponseData(body []byte) (*reqResp.StaticResponseConsumerAPI, error) {
	var resp reqResp.StaticResponseConsumerAPI
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
