package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aalu1418/rpi-terminal-hub/services/base"
	"github.com/aalu1418/rpi-terminal-hub/types"
	log "github.com/sirupsen/logrus"
)

type service struct {
	types.Service
	client *http.Client
	url    string
}

// TODO: remove, pass as flag
var KEY string

func init() {
	KEY = os.Getenv("OWM_KEY")
}

// provides a single handler for setting & incrementing metrics
func New(outgoingMsg chan<- types.Message) types.Service {
	var s service
	s.Service = base.New(outgoingMsg, types.WEATHER, types.WEATHER_FREQUENCY, s.onTick, s.processMsg)
	s.client = &http.Client{
		Timeout: types.DEFAULT_TIMEOUT,
	}
	s.url = queryBuilder(types.WEATHER_LON, types.WEATHER_LAT, KEY)
	return &s
}

func (s *service) processMsg(m types.Message) {
	log.Warnf("[%s] received unexpected message: %s", s.Name(), m)
}

// poll & send message to metrics service
func (s *service) onTick() (msg types.Message) {
	msg = types.Message{
		To:   types.WEBSERVER,
		Data: nil,
	}

	res, err := s.client.Get(s.url)
	if err != nil || res == nil || res.StatusCode != 200 {
		log.Errorf("failed to fetch weather: %+v, %s", res, err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Errorf("failed to read response body: %s", err)
	}
	var data types.OneCall
	if err := json.Unmarshal(body, &data); err != nil {
		log.Errorf("failed to unmarshal response body: %s", err)
	}

	fmt.Printf("%+v\n", data)

	return msg
}

func queryBuilder(lon, lat, key string) string {
	return fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%s&lon=%s&appid=%s&units=imperial&exclude=minutely,daily", lat, lon, key)
}