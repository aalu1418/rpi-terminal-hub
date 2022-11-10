package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

type service struct {
	types.Service
	client *http.Client
	url    string
}

// provides a single handler for setting & incrementing metrics
func New(outgoingMsg chan<- types.Message, key string) types.Service {
	var s service
	s.Service = base.New(outgoingMsg, types.WEATHER, types.WEATHER_FREQUENCY, s.onTick, s.processMsg)
	s.client = &http.Client{
		Timeout: types.DEFAULT_TIMEOUT,
	}
	s.url = queryBuilder(types.WEATHER_LON, types.WEATHER_LAT, key)
	return &s
}

func (s *service) processMsg(m types.Message) {
	log.Warnf("[%s] received unexpected message: %s", s.Name(), m)
}

// poll & send message to metrics service
func (s *service) onTick() (msg types.Message) {
	msg = types.Message{
		To: types.POSTOFFICE,
	}

	res, err := s.client.Get(s.url)
	if err != nil || res == nil || res.StatusCode != 200 {
		msg.Data = map[string]interface{}{
			"err": err,
			"res": res,
		}
		return msg
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		msg.Data = err
		return msg
	}
	var data types.OneCall
	if err := json.Unmarshal(body, &data); err != nil {
		msg.Data = err
		return msg
	}

	return types.Message{
		To:   types.WEBSERVER,
		Data: Parser(data),
	}
}

func queryBuilder(lon, lat, key string) string {
	return fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%s&lon=%s&appid=%s&units=imperial&exclude=minutely,daily", lat, lon, key)
}