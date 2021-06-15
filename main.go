package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

var apiKey string
var reg = prometheus.NewRegistry()
var client = &http.Client{Timeout: 10 * time.Second}

func initLog() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("logger initialised")
	apiKey = os.Args[2]
}

func init() {
	initLog()
	reg.MustRegister(NumDisconnections)
	reg.MustRegister(TotalUptimeSeconds)
	reg.MustRegister(TotalUptimePercentage)
	reg.MustRegister(AverageUptimeScore)
	reg.MustRegister(AverageLatencyMs)
	reg.MustRegister(AverageLatencyScore)
	reg.MustRegister(TotalRequests)
	reg.MustRegister(AverageSuccessRate)
	reg.MustRegister(AverageSuccessRateScore)
	reg.MustRegister(AverageReliabilityScore)
	reg.MustRegister(AverageUtilizationPercent)
	reg.MustRegister(TotalPreEarned)
	reg.MustRegister(CurrentlyConnected)
}

func childResult(metric string, c *gabs.Container) float64 {
	return c.S("period", metric).Data().(float64)
}

func updatePresearchMetric(metric string, c *gabs.Container, n string, g *prometheus.GaugeVec) {
	g.WithLabelValues(n).Set(childResult(metric, c))
}

func booltofloat64(inputbool bool) float64 {
	boolvar := float64(0)
	if inputbool {
		boolvar = 1
	}
	return boolvar
}

func childProcessor(children []*gabs.Container, node string) {
	c := children[0]

	currentlyConnected, _ := c.S("status", "connected").Data().(bool)
	CurrentlyConnected.WithLabelValues(node).Set(booltofloat64(currentlyConnected))

	disconnections, _ := c.S("period", "disconnections", "num_disconnections").Data().(float64)
	NumDisconnections.WithLabelValues(node).Set(disconnections)

	updatePresearchMetric("total_uptime_seconds", c, node, TotalUptimeSeconds)
	updatePresearchMetric("uptime_percentage", c, node, TotalUptimePercentage)
	updatePresearchMetric("avg_uptime_score", c, node, AverageUptimeScore)
	updatePresearchMetric("avg_latency_ms", c, node, AverageLatencyMs)
	updatePresearchMetric("avg_latency_score", c, node, AverageLatencyScore)
	updatePresearchMetric("total_requests", c, node, TotalRequests)
	updatePresearchMetric("avg_success_rate", c, node, AverageSuccessRate)
	updatePresearchMetric("avg_success_rate_score", c, node, AverageSuccessRateScore)
	updatePresearchMetric("avg_reliability_score", c, node, AverageReliabilityScore)
	updatePresearchMetric("avg_utilization_percent", c, node, AverageUtilizationPercent)
	updatePresearchMetric("avg_staked_capacity_percent", c, node, AverageStakedCapacityPercent)
	updatePresearchMetric("total_pre_earned", c, node, TotalPreEarned)
}

func checkNodeName(children []*gabs.Container) string {
	return children[0].S("meta", "description").Data().(string)
}

func presearchStatsHandler(w http.ResponseWriter, r *http.Request) {
	nodepublickey := r.FormValue("node")
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
	r.HandleFunc("/probe", presearchStatsHandler)

	log.Fatal(http.ListenAndServe(":"+os.Args[1], r))
}
