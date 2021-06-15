package main

import "github.com/prometheus/client_golang/prometheus"

var (
	NumDisconnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_number_of_disconnections",
			Help: "total number of disconnections",
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

	AverageLatencyMs = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "average_latency_ms",
			Help: "average latency in ms",
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

	AverageSuccessRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "average_success_rate",
			Help: "average success rate",
		},
		[]string{"nodename"},
	)

	AverageSuccessRateScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "average_success_rate_score",
			Help: "average success rate score",
		},
		[]string{"nodename"},
	)

	AverageReliabilityScore = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "average_reliability_score",
			Help: "average reliability score",
		},
		[]string{"nodename"},
	)

	AverageStakedCapacityPercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "avg_staked_capacity_percent",
			Help: "avg staked capacity percent",
		},
		[]string{"nodename"},
	)

	AverageUtilizationPercent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "avg_utilization_percent",
			Help: "average utilization percentage",
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

	CurrentlyConnected = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "currently_connected",
			Help: "A bool value for connection state of the node",
		},
		[]string{"nodename"},
	)
)
