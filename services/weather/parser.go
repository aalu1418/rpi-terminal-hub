package weather

import (
	"fmt"
	"strings"
	"time"

	"github.com/aalu1418/rpi-terminal-hub/types"
)

func Parser(data types.OneCall) types.WeatherParsed {
	// func for converting API timestamp to local time object
	convertTime := func(unix int64) time.Time {
		return time.Unix(unix, 0)
	}

	// func for converting Weather to WeatherParsedStatus
	convertWeather := func(w types.Weather) types.WeatherParsedStatus {
		out := types.WeatherParsedStatus{
			CloudCover: w.Clouds,
			UVIndex:    w.UVI,
			Humidity:   w.Humidity,
			Sunrise:    convertTime(w.Sunrise).Format("3:04 PM"),
			Sunset:     convertTime(w.Sunset).Format("3:04 PM"),
			Hour:       convertTime(w.Dt).Format("3 PM"),
			// Icon:         iconParser(w.Weather.Icon),
			Temperature:  int(w.Temp),
			FeelsLike:    int(w.FeelsLike),
			PrecipTotal:  0,
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
			out.Icon = "wi-na"
		}

		// wind speed formatting
		out.WindRange = fmt.Sprintf("%d", int(w.WindSpeed))
		if int(w.WindGust) != 0 && int(w.WindSpeed) != int(w.WindGust) {
			out.WindRange = out.WindRange + "-" + fmt.Sprintf("%d", int(w.WindGust))
		}

		return out
	}

	// set chance of precipitation on current
	data.Current.PrecipitationChance = data.Hourly[0].PrecipitationChance

	out := types.WeatherParsed{
		Date:        convertTime(data.Current.Dt).Format("Mon 02 Nov"),
		LastUpdated: time.Now().Format("Nov 2 @ 3:04 PM"),
		Current:     convertWeather(data.Current),
		Hourly:      [4]types.WeatherParsedStatus{},
	}

	// handle hourly
	for i := 0; i < 4; i++ {
		index := (1 + 3*i) % len(data.Hourly)
		out.Hourly[i] = convertWeather(data.Hourly[index])
	}

	return out
}

func iconParser(icon string) string {
	dayNight := "day"
	if icon[len(icon)-1] == byte('d') {
		dayNight = "night-alt"
	}

	iconMap := types.IconMap
	name, ok := iconMap[string(icon[:2])]
	if !ok {
		name = iconMap[icon]
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