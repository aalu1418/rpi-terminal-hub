package alerts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

type nwsService struct {
	types.Service
	client *http.Client
}

func NewNWS(outgoingMsg chan<- types.Message) types.Service {
	var s nwsService
	s.Service = base.New(outgoingMsg, types.NWS, types.NWS_FREQUENCY, s.onTick, s.processMsg)
	s.client = &http.Client{
		Timeout: types.DEFAULT_TIMEOUT,
	}
	return &s
}

func (s *nwsService) processMsg(m types.Message) {
	log.Warnf("[%s] received unexpected message: %s", s.Name(), m)
}

// poll & send message to metrics service
func (s *nwsService) onTick() types.Message {
	msg := types.Message{
		To: types.POSTOFFICE,
	}

	res, err := s.client.Get(fmt.Sprintf("https://api.weather.gov/alerts/active?point=%s,%s", types.WEATHER_LAT, types.WEATHER_LON))
	if err != nil {
		msg.Data = err
		return msg
	}
	if res == nil {
		msg.Data = fmt.Errorf("[CRITICAL] res is nil")
		return msg
	}
	if res.StatusCode != 200 {
		msg.Data = *res
		return msg
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		msg.Data = err
		return msg
	}

	var alerts types.NWSAlerts
	if err := json.Unmarshal(body, &alerts); err != nil {
		msg.Data = err
		return msg
	}

	// coalesce alerts
	out := map[string]types.ParsedAlert{}
	for _, v := range alerts.Features {
		severity := v.Properties.Severity
		if len(severity) > 3 {
			severity = severity[:3]
		}

		ind := strings.LastIndex(strings.TrimSpace(v.Properties.Event), " ")
		sev := types.ParseAlertLevel(v.Properties.Event[ind+1:])
		heading := fmt.Sprintf("[%s] %s",
			severity,
			v.Properties.Event[:ind],
		)

		if a, exists := out[heading]; exists {
			// handle increased level (overwrite)
			if sev > a.Level {
				a.Level = sev
				a.Start = v.Properties.Effective
				a.End = v.Properties.Ends

				out[heading] = a
				continue
			}

			// skip processing if lower level
			if sev < a.Level {
				continue
			}

			// if same level
			// handling combining times
			if v.Properties.Effective.Before(a.Start) {
				a.Start = v.Properties.Effective
			}
			if a.End.Before(v.Properties.Ends) {
				a.End = v.Properties.Ends
			}

			out[heading] = a
			continue
		}

		out[heading] = types.ParsedAlert{
			Level: sev,
			Start: v.Properties.Effective,
			End:   v.Properties.Ends,
		}
	}

	// printing
	output := []string{}
	if len(out) == 0 {
		output = append(output, "N/A")
	} else {
		for k, v := range out {
			// only show second date if different
			format := "1/02 3PM"
			if v.Start.Day() == v.End.Day() && v.Start.Month() == v.End.Month() {
				format = "3PM"
			}
			output = append(output, fmt.Sprintf(
				"%s %s: %s - %s",
				k,
				strings.Title(v.Level.String()),
				v.Start.Format("1/02 3PM"),
				v.End.Format(format),
			))
		}
	}

	return types.Message{
		To:   types.WEBSERVER,
		Data: strings.Join(output, ", "),
	}
}
