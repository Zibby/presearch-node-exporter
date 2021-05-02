# presearch prometheus node exporter

## Useage

```
./presearch-node-exporter $PORTTOBINDTO $YOUR_PRESEARCH_APIKEY
curl localhost:/probe/$NODE_PUBLIC_KEY # you will need to urlencode this key. eg use https://www.url-encode-decode.com/
```

API key and node public keys can be found at [https://nodes.presearch.org/dashboard]()

## Example output

```
# HELP average_latency_ms average latency in ms
# TYPE average_latency_ms gauge
average_latency_ms{nodename="Linode"} 1018.3943
# HELP average_latency_score Average latency score
# TYPE average_latency_score gauge
average_latency_score{nodename="Linode"} 27.180229
# HELP average_reliability_score average reliability score
# TYPE average_reliability_score gauge
average_reliability_score{nodename="Linode"} 75.36464
# HELP average_success_rate average success rate
# TYPE average_success_rate gauge
average_success_rate{nodename="Linode"} 100
# HELP average_success_rate_score average success rate score
# TYPE average_success_rate_score gauge
average_success_rate_score{nodename="Linode"} 100
# HELP average_uptime_score Average uptime score
# TYPE average_uptime_score gauge
average_uptime_score{nodename="Linode"} 99.857763
# HELP avg_utilization_percent average utilization percentage
# TYPE avg_utilization_percent gauge
avg_utilization_percent{nodename="Linode"} 1.01717e-05
# HELP total_number_of_disconnections total number of disconnections
# TYPE total_number_of_disconnections gauge
total_number_of_disconnections{nodename="Linode"} 25
# HELP total_number_of_requests total number of requests
# TYPE total_number_of_requests gauge
total_number_of_requests{nodename="Linode"} 175
# HELP total_pre_earned Total Pre Earned
# TYPE total_pre_earned gauge
total_pre_earned{nodename="Linode"} 28.69725493086953
# HELP total_uptime_percentage total uptime percentage
# TYPE total_uptime_percentage gauge
total_uptime_percentage{nodename="Linode"} 99.59968435222578
# HELP total_uptime_seconds Total uptime of node in seconds
# TYPE total_uptime_seconds gauge
total_uptime_seconds{nodename="Linode"} 2.145676e+06
```