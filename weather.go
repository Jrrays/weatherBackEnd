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

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	}
}

func (w *WeatherResponse) String() string {
	return fmt.Sprintf("City: %s\nTemp: %.1f °C\nHumidity: %d %%\nConditoins: %s\n",
		w.Name,
		w.Main.Temp,
		w.Main.Humidity,
		w.Weather[0].Description)
}

type ClientResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature_celsius"`
	Humidity    int     `json:"humidity_percent"`
	Condition   string  `json:"description"`
}

func UrlBuilder() (string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "api.openweathermap.org",
	}

	u = u.JoinPath("data", "2.5", "weather")

	q := u.Query()
	q.Set("lat", os.Getenv("LAT"))
	q.Set("lon", os.Getenv("LON"))
	q.Set("appid", os.Getenv("API"))
	q.Set("units", os.Getenv("UNITS"))
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func fetchWeatherData() (*WeatherResponse, error) {
	//Build the URL
	apiURL, err := UrlBuilder()
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

func currentweather(w http.ResponseWriter, req *http.Request) {

	//Check if the req method == GET
	if req.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	//Fetch upstream data
	upstreamData, err := fetchWeatherData()
	if err != nil {
		log.Println("Error fetching weather:", err)
		http.Error(w, "Failed to retrieve weather data", http.StatusInternalServerError)
		return
	}

	//Format the weather data
	cond := "unknown"
	if len(upstreamData.Weather) > 0 {
		cond = upstreamData.Weather[0].Description
	}

	clientData := ClientResponse{
		City:        upstreamData.Name,
		Temperature: upstreamData.Main.Temp,
		Humidity:    upstreamData.Main.Humidity,
		Condition:   cond,
	}

	//Set headers and return JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(clientData); err != nil {
		log.Println("Error while encoding response:", err)
	}
}

func main() {

	port := os.Getenv("PORT")
	addr := ":" + port

	http.HandleFunc("/currentweather", currentweather)

	fmt.Println("Server listening port", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}
