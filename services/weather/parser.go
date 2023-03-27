package weather

import (
	"fmt"
	"strings"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/types"
)

func Parser(data types.OneCall) types.WeatherParsed {
	// set chance of precipitation on current
	if len(data.Hourly) > 0 {
		data.Current.PrecipitationChance = data.Hourly[0].PrecipitationChance
	}

	out := types.WeatherParsed{
		Date:        timeParser(data.Current.Dt).Format("Mon 02 Jan"),
		LastUpdated: time.Now().Format("Jan 2 @ 3:04 PM"),
		Location:    types.WEATHER_LOCATION,
		Current:     weatherParser(data.Current),
		Hourly:      [4]types.WeatherParsedStatus{},
	}

	// handle hourly
	for i := 0; i < 4; i++ {
		if len(data.Hourly) > 0 {
			index := (1 + 3*i) % len(data.Hourly)
			out.Hourly[i] = weatherParser(data.Hourly[index])
		} else {
			out.Hourly[i] = weatherParser(data.Current)
		}
	}

	return out
}

// func for converting API timestamp to local time object
func timeParser(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// func for converting Weather to WeatherParsedStatus
func weatherParser(w types.Weather) types.WeatherParsedStatus {
	out := types.WeatherParsedStatus{
		CloudCover: w.Clouds,
		UVIndex:    w.UVI,
		Humidity:   w.Humidity,
		Sunrise:    timeParser(w.Sunrise).Format("3:04 PM"),
		Sunset:     timeParser(w.Sunset).Format("3:04 PM"),
		Hour:       timeParser(w.Dt).Format("3 PM"),
		// Icon:         iconParser(w.Weather.Icon),
		Temperature:  int(w.Temp),
		FeelsLike:    fmt.Sprintf("%.1f", w.FeelsLike),
		PrecipTotal:  fmt.Sprintf("%.1f", w.Rain.OneHour+w.Snow.OneHour),
		PrecipChance: int(w.PrecipitationChance * 100),
		// WindRange:    "", // set below
	}

	// hour padding
	if len(out.Hour) == 4 {
		out.Hour = " " + out.Hour
	}

	// Icon handling
	if len(w.Weather) > 0 {
		out.Icon = iconParser(w.Weather[0].Icon)
	} else {
		out.Icon = types.IconNA
	}

	// wind speed formatting
	out.WindRange = fmt.Sprintf("%d", int(w.WindSpeed))
	if int(w.WindGust) != 0 && int(w.WindSpeed) < int(w.WindGust) {
		out.WindRange = out.WindRange + "-" + fmt.Sprintf("%d", int(w.WindGust))
	}

	return out
}

func iconParser(icon string) string {
	if len(icon) != 3 {
		return types.IconNA
	}

	dayNight := "day"
	if icon[len(icon)-1] != byte('d') {
		dayNight = "night-alt"
	}

	iconMap := types.IconMap
	name, ok := iconMap[string(icon[:2])]
	if !ok {
		name, ok = iconMap[icon]
		if !ok {
			name = types.IconNA // default if not found
		}
	}

	if strings.Contains(name, "wi-") {
		return name
	}

	return strings.Join([]string{
		types.IconPrefix,
		dayNight,
		name,
	}, "-")
}