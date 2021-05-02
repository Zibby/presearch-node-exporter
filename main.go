package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

var client = &http.Client{Timeout: 10 * time.Second}
var apiKey string

var (
	reg = prometheus.NewRegistry()

	TotalUptimePercentage = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_uptime_percentage",
			Help: "total uptime percentage",
		},
	)

	AverageUptimeScore = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "average_uptime_score",
			Help: "Average uptime score",
		},
	)

	AverageLatencyScore = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "average_latency_score",
			Help: "Average latency score",
		},
	)

	TotalUptimeSeconds = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_uptime_seconds",
			Help: "Total uptime of node in seconds",
		},
	)

	TotalPreEarned = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_pre_earned",
			Help: "Total Pre Earned",
		},
	)
)

func initLog() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("logger initialised")
	apiKey = os.Args[1]
}

func init() {
	initLog()
	reg.MustRegister(TotalUptimeSeconds)
	reg.MustRegister(AverageLatencyScore)
	reg.MustRegister(AverageUptimeScore)
	reg.MustRegister(TotalUptimePercentage)
	reg.MustRegister(TotalPreEarned)
}

func semiescapeurl(i string) string {
	i = strings.ReplaceAll(i, "%2F", `/`)
	i = strings.ReplaceAll(i, "%5C", `\`)
	i = strings.ReplaceAll(i, "%2B", `+`)
	i = strings.ReplaceAll(i, "%0D", `\n`)
	i = strings.ReplaceAll(i, "%0A", "")
	return i
}

func childProcessor(children []*gabs.Container) {
	ut, _ := children[0].S("period", "total_uptime_seconds").Data().(float64)
	TotalUptimeSeconds.Set(ut)

	utp, _ := children[0].S("period", "uptime_percentage").Data().(float64)
	TotalUptimePercentage.Set(utp)

	auts, _ := children[0].S("period", "avg_uptime_score").Data().(float64)
	AverageUptimeScore.Set(auts)

	als, _ := children[0].S("period", "average_latency_score").Data().(float64)
	AverageLatencyScore.Set(als)

	te, _ := children[0].S("period", "total_pre_earned").Data().(float64)
	AverageLatencyScore.Set(te)
}

func presearchStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var nodepublickey string
	nodepublickey = vars["node"]
	resp, err := client.Get("https://nodes.presearch.org/api/nodes/status/" + apiKey + "?stats=true&nodes=" + nodepublickey)
	if err != nil {
		log.Error("Failed to connect to api")
	}
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read body")
	}
	jsonParsed, err := gabs.ParseJSON(ret)
	defer resp.Body.Close()

	if err != nil {
		log.Error("Failed to decode node name")
	}
	children := jsonParsed.Path("nodes").Children()
	childProcessor(children)
	promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
}

func healthHander(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Health Ok!")
}

func main() {
	r := mux.NewRouter()
	r.UseEncodedPath()
	r.HandleFunc("/health", healthHander)
	r.HandleFunc("/probe/{node}", presearchStatsHandler)

	log.Fatal(http.ListenAndServe(":8082", r))
}
