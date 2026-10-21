package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

// Data from OpenWeatherMap is stored in this struct
type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		WindSpeed float64 `json:"speed"`
	} `json:"wind"`
	Weather []struct {
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}
}

// Struct for client response
type ClientResponse struct {
	City        string  `json:"City"`
	Temperature float64 `json:"Temperature"`
	FeelsLike   float64 `json:"FeelsLike"`
	Humidity    int     `json:"Humidity"`
	WindSpeed   float64 `json:"WindSpeed"`
	Condition   string  `json:"Description"`
	Icon        string  `json:"Icon"`
}

// String function for WeatherResponse to client
func (w *WeatherResponse) String() string {
	return fmt.Sprintf("City: %s\nTemp: %.1f °C\nFeels Like: %.1f°C\nHumidity: %d %%\nWindSpeed: %f m/s\nConditoins: %s\nIcon: %s\n",
		w.Name,
		w.Main.Temp,
		w.Main.FeelsLike,
		w.Main.Humidity,
		w.Wind.WindSpeed,
		w.Weather[0].Description,
		w.Weather[0].Icon,
	)
}

func UrlBuilder(lat, lon string) (string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "api.openweathermap.org",
	}

	u = u.JoinPath("data", "2.5", "weather")

	q := u.Query()
	q.Set("lat", lat)
	q.Set("lon", lon)
	q.Set("appid", os.Getenv("API"))
	q.Set("units", os.Getenv("UNITS"))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func fetchWeatherData(lat, lon string) (*WeatherResponse, error) {
	//Build the URL
	apiURL, err := UrlBuilder(lat, lon)
	if err != nil {
		return nil, err
	}

	//Fetch the weather data
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to Reach OpenWeatherMap: %w", err)
	}
	defer resp.Body.Close()

	//Check whether the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %s", resp.Status)
	}

	var data WeatherResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Failed to parse upstream JSON: %w", err)
	}

	return &data, nil
}

// Handle requests from the client
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
	upstreamData, err := fetchWeatherData(lat, lon)
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

	clientData := ClientResponse{
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
