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

	out := []string{}
	for _, v := range alerts.Features {
		// only show second date if different
		format := "1/02 3PM"
		if v.Properties.Effective.Day() == v.Properties.Expires.Day() && v.Properties.Effective.Month() == v.Properties.Expires.Month() {
			format = "3PM"
		}

		severity := v.Properties.Severity
		if len(severity) > 3 {
			severity = severity[:3]
		}

		out = append(out, fmt.Sprintf(
			"[%s] %s: %s - %s",
			severity,
			v.Properties.Event,
			v.Properties.Effective.Format("1/02 3PM"),
			v.Properties.Expires.Format(format),
		))
	}

	if len(out) == 0 {
		out = append(out, "N/A")
	}

	return types.Message{
		To:   types.WEBSERVER,
		Data: strings.Join(out, ", "),
	}
}
