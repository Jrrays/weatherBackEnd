package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"weatherBackEnd/weatherutils"

	_ "github.com/joho/godotenv/autoload"
)

// Handler for weather requests
func currentweather(w http.ResponseWriter, req *http.Request) {

	//Check if the req method == GET
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Extract location from query parameters
	params := req.URL.Query()
	lat := params.Get("lat")
	lon := params.Get("lon")

	//Fetch upstream data
	upstreamData, err := weatherutils.FetchWeatherData(lat, lon)
	if err != nil {
		log.Println("Error fetching weather:", err)
		http.Error(w, "Failed to retrieve weather data", http.StatusInternalServerError)
		return
	}

	//Format the weather data
	cond := "unknown"
	icon := "unknown"
	if len(upstreamData.Weather) > 0 {
		cond = upstreamData.Weather[0].Description
		icon = upstreamData.Weather[0].Icon
	}

	clientData := weatherutils.ClientResponse{
		City:        upstreamData.Name,
		Temperature: upstreamData.Main.Temp,
		FeelsLike:   upstreamData.Main.FeelsLike,
		Humidity:    upstreamData.Main.Humidity,
		WindSpeed:   upstreamData.Wind.WindSpeed,
		Condition:   cond,
		Icon:        icon,
	}

	//Set headers and return JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(clientData); err != nil {
		log.Println("Error while encoding response:", err)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Hardcoded Access-Control-Allow-Origin value. Maybe fix this later.
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	addr := ":" + port

	mux := http.NewServeMux()
	mux.HandleFunc("/currentweather", currentweather)

	fmt.Println("Server listening port", port)
	if err := http.ListenAndServe(addr, enableCORS(mux)); err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}
