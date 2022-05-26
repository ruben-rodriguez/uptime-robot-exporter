package main
          
import (
    "strings"
    "net/http"
    "encoding/json"
    "os"
	"time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Monitor struct {
    ID              int         `json:"id"`
	FriendlyName    string      `json:"friendly_name"`
	URL             string      `json:"url"`
	Type            string      `json:"type"`
	SubType         string      `json:"sub_type"`
	KeywordType     interface{} `json:"keyword_type"`
	KeywordCaseType interface{} `json:"keyword_case_type"`
	KeywordValue    string      `json:"keyword_value"`
	HTTPUsername    string      `json:"http_username"`
	HTTPPassword    string      `json:"http_password"`
	Port            string      `json:"port"`
	Interval        int         `json:"interval"`
	Timeout         int         `json:"timeout"`
	Status          int         `json:"status"`
	CreateDatetime  int         `json:"create_datetime"`
}

type MonitorResponse struct {
    Data []Monitor `json:"monitors"`
}

const (
	url = "https://api.uptimerobot.com/v2/getMonitors"
)

var (
    uptimeRobotAPIKey       =  getEnv("UPTIME_ROBOT_API_KEY")
    pollingInterval         =  getEnv("POLLING_INTERVAL")

    status = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "uptime_robot_status",
            Help: "Uptime Robot status",
		},
		[]string{"name", "url", "type", "sub_type", "port"},
	)
)

func main() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    zerolog.SetGlobalLevel(zerolog.DebugLevel)

    interval, err := time.ParseDuration(pollingInterval)
    if err != nil {
        log.Error().Err(err)
        panic(err)
    }

    log.Info().Msgf("Polling interval %s", interval)

    http.Handle("/metrics", promhttp.Handler())
    go http.ListenAndServe(":9101", nil)

    for _ = range time.Tick(interval) {
		timestamp := time.Now()

		str := "Polling Uptime Robot API data at " + timestamp.String()
		log.Info().Msg(str)

        for _, element := range getMonitors() {
            status.WithLabelValues(element.FriendlyName, element.URL, element.Type, 
            element.SubType, element.Port).Set(float64(element.Status))

        }

	} 

}

func getEnv(key string) string {
    value, ok := os.LookupEnv(key)
    if !ok {

        log.Fatal().Msgf("%s ENV var is not present", key)
        os.Exit(1)
    } 
    return value
}

func getMonitors() []Monitor {

    var response MonitorResponse
    var monitors []Monitor
    
    payload := strings.NewReader("api_key=" + uptimeRobotAPIKey + "&format=json&logs=0")
        
    req, _ := http.NewRequest("POST", url, payload)
          
    req.Header.Add("content-type", "application/x-www-form-urlencoded")
    req.Header.Add("cache-control", "no-cache")
          
    res, err := http.DefaultClient.Do(req)
    if err != nil {
		log.Error().Err(err)
	}

    defer res.Body.Close()
    if res.StatusCode != 200 {
		log.Info().Int("status", res.StatusCode).Msg("Request status code")
	}

    if res.StatusCode == http.StatusOK {

        json.NewDecoder(res.Body).Decode(&response)

        if len(response.Data) > 0 {

            log.Info().Int("number of monitors",len(response.Data))

            for _, monitor := range response.Data {

                monitors = append(monitors, Monitor {
                    ID:                 monitor.ID,
                    FriendlyName:       monitor.FriendlyName,
                    URL:                monitor.URL,
                    Type:               monitor.Type,
                    Port:               monitor.Port,
                    Interval:           monitor.Interval,
                    Status:             monitor.Status,
                    CreateDatetime:     monitor.CreateDatetime,
                })

            }

        } else {
            log.Warn().Msg("Received empty response from Uptime Robot")

        }

    }

    return monitors
}
