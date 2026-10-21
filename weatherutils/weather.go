package weatherutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Data from OpenWeatherMap is stored in this struct
type weatherResponse struct {
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
func (w *weatherResponse) String() string {
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

func urlBuilder(lat, lon string) (string, error) {
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

func FetchWeatherData(lat, lon string) (*weatherResponse, error) {
	//Build the URL
	apiURL, err := urlBuilder(lat, lon)
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

	var data weatherResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Failed to parse upstream JSON: %w", err)
	}

	return &data, nil
}
