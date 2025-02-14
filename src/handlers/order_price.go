package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/your-username/your-project/src/constants"
	custErr "github.com/your-username/your-project/src/errors"
	service "github.com/your-username/your-project/src/services"
	// "github.com/your-username/your-project/internal/services"
)

type OrderPriceHandler struct {
	calculator service.DeliveryPriceCalculator
}

// NewOrderPriceHandler is a constructor for the handler
// this can be done for all the individual handlers that you will be having for each of the service handing a logical task (Separation of concerns.)
func NewOrderPriceHandler(c service.DeliveryPriceCalculator) http.HandlerFunc {
	h := &OrderPriceHandler{calculator: c}
	return h.deliveryOrderPriceHandler
}

// DeliveryOrderPriceHandler handles the GET /api/v1/delivery-order-price request.
// the rest of the endpoints pertaining the order can be placed in this file and exposed with the above handler
func (h *OrderPriceHandler) deliveryOrderPriceHandler(w http.ResponseWriter, r *http.Request) {
	// Read query parameters
	venueSlug := r.URL.Query().Get("venue_slug")
	cartValueStr := r.URL.Query().Get("cart_value")
	userLatStr := r.URL.Query().Get("user_lat")
	userLonStr := r.URL.Query().Get("user_lon")

	if venueSlug == "" || cartValueStr == "" || userLatStr == "" || userLonStr == "" {
		errorString := fmt.Sprintf("Missing required query parameters. Please check your request. Statuscode %v", http.StatusBadRequest)
		log.Println(errorString)
		http.Error(w, errorString, http.StatusBadRequest)
		return
	}

	cartValue, err := strconv.Atoi(cartValueStr)
	if err != nil {
		log.Printf("Invalid cartValue in queryparams:")
		http.Error(w, "Invalid cart_value", http.StatusBadRequest)
		return
	}

	userLat, err := strconv.ParseFloat(userLatStr, 64)
	if err != nil {
		log.Printf("Invalid latitude in queryparams: %v", err)
		http.Error(w, "Invalid user_lat", http.StatusBadRequest)
		return
	}

	userLon, err := strconv.ParseFloat(userLonStr, 64)
	if err != nil {
		log.Printf("Invalid longitude in queryparams: %v", err)
		http.Error(w, "Invalid user_lon", http.StatusBadRequest)
		return
	}

	log.Printf("Calculating delivery price for venue_slug=%s cart_value=%d user_lat=%f user_lon=%f",
		venueSlug, cartValue, userLat, userLon) //add in globalCustomerID here for better debugging of the

	// Calculate the order price
	responseData, err := h.calculator.CalculateDeliveryOrderPrice(venueSlug, cartValue, userLat, userLon)
	if err != nil {
		var svcErr *custErr.ServiceError
		if errors.As(err, &svcErr) {
			if svcErr.Code == constants.LOCATION_SERVICE_DISTANCE_OUT_OF_REACH_400 {
				// Handle the specific error case here.
				log.Printf("distance is out of reach " + svcErr.Code)
				http.Error(w, "distance is out of reach", http.StatusBadRequest)
				return
			}
		}

		// Otherwise, respond with a generic 500
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
		// errorString := fmt.Sprintf("Error in calculating Delivery Order Price: %v", err)
		// log.Print(errorString)
		// http.Error(w, "Internal Server Error", http.StatusBadRequest)
		// return
	}
	// Marshal the response to JSON and write to output
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		log.Printf("Error in encoding response data for the json response: %v", err)
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
		return
	}
}
