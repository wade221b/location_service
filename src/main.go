package main

import (
	"log"
	"net/http"

	httpclient "github.com/your-username/your-project/src/client"
	"github.com/your-username/your-project/src/config"
	"github.com/your-username/your-project/src/handlers"
	service "github.com/your-username/your-project/src/services"
)

func main() {
	// Load configuration (environment variables, etc.)
	cfg := config.LoadConfig()

	// Initialize the HTTP router (in this example, we just use http.DefaultServeMux)
	//for each new handler, pass the same http client of make new one for each new handler that you create.
	mux := http.NewServeMux()
	httpClient := httpclient.NewHTTPClient(cfg.ConsumerHost)
	priceCalculator := service.NewDeliveryPriceCalculator(httpClient)
	//as the httpClient is passed here, it can be done for all the respective api handlers that would be made.

	// 3. Build the delivery price calculator with the http client

	// Register your handlers
	mux.HandleFunc("/api/v1/delivery-order-price", handlers.NewOrderPriceHandler(priceCalculator))
	//the above mapping of adding handlers can be done in another file, wherein all the APIs will be listed.
	//here as only 1 api is present, i am keeping it openly visible in the main file.

	// Start the server
	addr := cfg.ServerAddress
	log.Printf("Starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
