package types

type OneCall struct {
	Lat            float64   `json:"lat"`
	Lon            float64   `json:"lon"`
	Timezone       string    `json:"timezone"`
	TimezoneOffset int64     `json:"timezone_offset"`
	Current        Weather   `json:"current"`
	Hourly         []Weather `json:"hourly"`
}

type Weather struct {
	Dt                  int64         `json:"dt"`
	Sunrise             int64         `json:"sunrise"` // only on current
	Sunset              int64         `json:"sunset"`  // only on current
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
	OneHour float64 `json:"1h,omitempty"`
}

type Details struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type WeatherParsed struct {
	Date        string
	LastUpdated string
	Current     WeatherParsedStatus
	Hourly      [4]WeatherParsedStatus
}

type WeatherParsedStatus struct {
	// parameters only needed for Current
	CloudCover int
	UVIndex    float64
	Humidity   int
	Sunrise    string
	Sunset     string

	// parameters needed for both
	Hour         string
	Icon         string
	Temperature  int
	FeelsLike    int
	PrecipTotal  int
	PrecipChance int
	WindRange    string
}

const IconPrefix = "wi"

var IconMap = map[string]string{
	"01d": "sunny",
	"01n": "wi-night-clear",
	"02d": "sunny-overcast",
	"02n": "partly-cloudy",
	"03":  "cloudy",
	"04":  "wi-cloudy",
	"09":  "showers",
	"10":  "rain",
	"11":  "thunderstorm",
	"13":  "snow",
	"50d": "fog",
	"50n": "wi-night-fog",
}