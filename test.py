#!/usr/bin/env python3

import requests

def main():
    base_url = "http://localhost:8000/api/v1/delivery-order-price"

    # Define test scenarios
    test_scenarios = [
        {
            "description": "Happy path - standard Helsinki values",
            "params": {
                "venue_slug": "home-assignment-venue-helsinki",
                "cart_value": "1000",
                "user_lat": "60.17094",
                "user_lon": "24.93087"
            }
        },
        {
            "description": "Stockholm venue, valid cart value",
            "params": {
                "venue_slug": "home-assignment-venue-stockholm",
                "cart_value": "1500",
                "user_lat": "59.3293",
                "user_lon": "18.0686"
            }
        },
        {
            "description": "Tokyo venue, large cart value, standard lat/lon",
            "params": {
                "venue_slug": "home-assignment-venue-tokyo",
                "cart_value": "5000",
                "user_lat": "35.6895",
                "user_lon": "139.6917"
            }
        },
        {
            "description": "Berlin venue, small cart value, standard lat/lon",
            "params": {
                "venue_slug": "home-assignment-venue-berlin",
                "cart_value": "500",
                "user_lat": "52.5200",
                "user_lon": "13.4050"
            }
        },
        {
            "description": "Missing user_lat to simulate an error scenario",
            "params": {
                "venue_slug": "home-assignment-venue",
                "cart_value": "1000",
                # intentionally missing user_lat
                "user_lon": "24.93087"
            }
        },
        {
            "description": "Missing venue_slug to simulate an error scenario",
            "params": {
                # intentionally missing venue_slug
                "cart_value": "1000",
                "user_lat": "60.17094",
                "user_lon": "24.93087"
            }
        },
        {
            "description": "Very large cart value, valid lat/lon",
            "params": {
                "venue_slug": "home-assignment-venue",
                "cart_value": "999999",
                "user_lat": "60.17094",
                "user_lon": "24.93087"
            }
        },
        {
            "description": "Very large coordinates to simulate out-of-range distance",
            "params": {
                "venue_slug": "home-assignment-venue",
                "cart_value": "1000",
                "user_lat": "85.0000",  # near the north pole
                "user_lon": "179.9999"
            }
        },
        {
            "description": "Empty cart_value",
            "params": {
                "venue_slug": "home-assignment-venue-helsinki",
                "cart_value": "",  # empty
                "user_lat": "60.17094",
                "user_lon": "24.93087"
            }
        }
    ]

    for scenario in test_scenarios:
        print(f"=== Test: {scenario['description']} ===")
        try:
            resp = requests.get(base_url, params=scenario["params"])
            
            print("HTTP Status Code:", resp.status_code)
            
            print("Response Body:", resp.text)
            
        except requests.exceptions.RequestException as e:
            print("Request failed:", e)
        
        print("\n------------------------------\n")

if __name__ == "__main__":
    main()
