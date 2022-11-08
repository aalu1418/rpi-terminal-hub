package types

type OneCall struct {
	Lat            float64   `json:"lat"`
	Lon            float64   `json:"lon"`
	Timezone       string    `json:"timezone"`
	TimezoneOffset int       `json:"timezone_offset"`
	Current        Weather   `json:"current"`
	Hourly         []Weather `json:"hourly"`
}

type Weather struct {
	Dt                  int           `json:"dt"`
	Sunrise             int           `json:"sunrise"` // only on current
	Sunset              int           `json:"sunset"`  // only on current
	Temp                float64       `json:"temp"`
	FeelsLike           float64       `json:"feels_like"`
	Pressure            int           `json:"pressure"`
	Humidity            int           `json:"humidity"`
	DewPoint            float64       `json:"dew_point"`
	Clouds              int           `json:"clouds"`
	UVI                 float64       `json:"uvi"`
	Visibility          int           `json:"visibility"`
	WindSpeed           float64       `json:"wind_speed"`
	WindGust            float64       `json:"wind_gust"`
	WindDeg             float64       `json:"wind_deg"`
	PrecipitationChance float64       `json:"pop"` // only on hourly
	Rain                Precipitation `json:"rain"`
	Snow                Precipitation `json:"snow"`
	Weather             []Details     `json:"weather"`
}

type Precipitation struct {
	OneHour   float64 `json:"1h,omitempty"`
	ThreeHour float64 `json:"3h,omitempty"`
}

type Details struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
