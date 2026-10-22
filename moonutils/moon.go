package moonutils

import "time"

type moonResponse struct {
	Name  string `json:"name"`
	Phase struct {
		PhaseName string  `json:"phase_name"`
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
	PhaseName string    `json:"phase_name"`
	Angle     float64   `json:"phase_angle_deg"`
	Age       float64   `json:"age_days"`
	Distance  float64   `json:"distance_km"`
	Svg       string    `json:"svg"`
	Name      string    `json:"name"`
	FullMoon  time.Time `json:"applies_to_full_moon_at"`
}
