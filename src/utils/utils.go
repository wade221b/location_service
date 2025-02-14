package utils

import (
	"log"
	"math"

	reqResp "github.com/your-username/your-project/src/req_resp"
)

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 { //what i first thought. basic distance between two points
	sol := math.Sqrt(math.Pow(lat1-lat2, 2) + math.Pow(lon1-lon2, 2))
	return sol
}

func CalculateExactDistance(lat1, lon1, lat2, lon2 float64) float64 { //reference : https://stackoverflow.com/a/27943/22299413
	const earthRadius = 6371000 // Earth radius in meters

	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	lat1Rad := lat1 * math.Pi / 180.0
	lat2Rad := lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1Rad)*math.Cos(lat2Rad)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return distance
}

func GetDistanceRange(dynamicVenueDetails *reqResp.DynamicResponseConsumerAPI, distance float64) *reqResp.DistanceRange {
	var distanceRanges []reqResp.DistanceRange = dynamicVenueDetails.VenueRaw.DeliverySpecs.DeliveryPricing.DistanceRanges

	log.Printf("distnce is %v", distance)
	log.Println("last values min value is ", distanceRanges[len(distanceRanges)-1].Min)

	if len(distanceRanges) == 0 {
		return nil
	}

	if distance >= float64(distanceRanges[len(distanceRanges)-1].Min) { //check with the last element if the delivery is possible
		return nil //check how this has to be handled
	}

	for _, distrange := range distanceRanges {
		if distance >= float64(distrange.Min) && distance < float64(distrange.Max) {
			return &distrange
		}
	}
	return nil
}
