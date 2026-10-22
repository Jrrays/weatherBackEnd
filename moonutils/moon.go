package moonutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type moonResponse struct {
	Name  string `json:"name"`
	Phase struct {
		PhaseName string  `json:"name"`
		Angle     float64 `json:"phase_angle_deg"`
		Age       float64 `json:"age_days"`
		Distance  float64 `json:"distance_km"`
	} `json:"phase"`
	Visual struct {
		Svg string `json:"svg"`
	} `json:"visual"`
	Traditional struct {
		Name     string    `json:"name"`
		FullMoon time.Time `json:"applies_to_full_moon_at"`
	} `json:"traditional_moon"`
}

type ClientResponse struct {
	PhaseName string    `json:"PhaseName"`
	Angle     float64   `json:"PhaseAngle"`
	Age       float64   `json:"Age"`
	Distance  float64   `json:"Distance"`
	Svg       string    `json:"Svg"`
	Name      string    `json:"TraditionalName"`
	FullMoon  time.Time `json:"AppliesTo"`
}

func urlBuilder(lat, lon string) (string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   "api.freeastroapi.com",
	}

	u = u.JoinPath("api", "v1", "moon", "phase")

	q := u.Query()
	q.Set("lat", lat)
	q.Set("lon", lon)
	q.Set("include_visuals", "true")
	q.Set("include_traditional_moon", "true")
	q.Set("include_forecast", "true")
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func FetchMoonData(lat, lon string) (*moonResponse, error) {
	//Form the query URL
	apiURL, err := urlBuilder(lat, lon)
	if err != nil {
		return nil, err
	}

	//Fetch the moon data
	client := http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("x-api-key", os.Getenv("FAA_API"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to reach FreeAstroAPI: %w", err)
	}
	defer resp.Body.Close()

	//Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %s", resp.Status)
	}

	var data moonResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Failed to parse upstream JSON: %w", err)
	}

	return &data, nil
}
