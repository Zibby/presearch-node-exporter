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

	TotalUptimePercentage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_uptime_percentage",
			Help: "total uptime percentage",
		},
		[]string{"nodename"},
	)

	AverageUptimeScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "average_uptime_score",
			Help: "Average uptime score",
		},
		[]string{"nodename"},
	)

	AverageLatencyScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "average_latency_score",
			Help: "Average latency score",
		},
		[]string{"nodename"},
	)

	TotalUptimeSeconds = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_uptime_seconds",
			Help: "Total uptime of node in seconds",
		},
		[]string{"nodename"},
	)

	TotalPreEarned = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_pre_earned",
			Help: "Total Pre Earned",
		},
		[]string{"nodename"},
	)

	NumDisconnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_number_of_disconnections",
			Help: "total number of disconnections",
		},
		[]string{"nodename"},
	)

	TotalRequests = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_number_of_requests",
			Help: "total number of requests",
		},
		[]string{"nodename"},
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
	reg.MustRegister(NumDisconnections)
	reg.MustRegister(TotalRequests)

}

func semiescapeurl(i string) string {
	i = strings.ReplaceAll(i, "%2F", `/`)
	i = strings.ReplaceAll(i, "%5C", `\`)
	i = strings.ReplaceAll(i, "%2B", `+`)
	i = strings.ReplaceAll(i, "%0D", `\n`)
	i = strings.ReplaceAll(i, "%0A", "")
	return i
}

func childProcessor(children []*gabs.Container, nodename string) {
	ut, _ := children[0].S("period", "total_uptime_seconds").Data().(float64)
	TotalUptimeSeconds.WithLabelValues(nodename).Set(ut)

	utp, _ := children[0].S("period", "uptime_percentage").Data().(float64)
	TotalUptimePercentage.WithLabelValues(nodename).Set(utp)

	auts, _ := children[0].S("period", "avg_uptime_score").Data().(float64)
	AverageUptimeScore.WithLabelValues(nodename).Set(auts)

	als, _ := children[0].S("period", "average_latency_score").Data().(float64)
	AverageLatencyScore.WithLabelValues(nodename).Set(als)

	te, _ := children[0].S("period", "total_pre_earned").Data().(float64)
	AverageLatencyScore.WithLabelValues(nodename).Set(te)

	disconnections, _ := children[0].S("period", "disconnections", "num_disconnections").Data().(float64)
	NumDisconnections.WithLabelValues(nodename).Set(disconnections)

  tr, _ := children[0].S("period", "total_requests").Data().(float64)
	TotalRequests.WithLabelValues(nodename).Set(tr)
}

func checkNodeName(children []*gabs.Container) string {
  return children[0].S("meta", "description").Data().(string)
}

func presearchStatsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var nodepublickey string
	nodepublickey = vars["node"]
	resp, err := client.Get("https://nodes.presearch.org/api/nodes/status/" + apiKey + "?stats=true&start_date=2001-01-01-00%3A00&nodes=" + nodepublickey)
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
	nodename := checkNodeName(children)
	childProcessor(children, nodename)
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
