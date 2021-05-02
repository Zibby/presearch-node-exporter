# presearch prometheus node exporter

## Useage

```
./presearch-node-exporter $YOUR_PRESEARCH_APIKEY
curl localhost:/probe/$NODE_PUBLIC_KEY # you will need to urlencode this key. eg use https://www.url-encode-decode.com/
```

API key and node public keys can be found at [https://nodes.presearch.org/dashboard]()

## Example output

```
# HELP average_latency_score Average latency score
# TYPE average_latency_score gauge
average_latency_score{nodename="Linode"} 27.921910365269877
# HELP average_uptime_score Average uptime score
# TYPE average_uptime_score gauge
average_uptime_score{nodename="Linode"} 99.856192
# HELP total_number_of_disconnections total number of disconnections
# TYPE total_number_of_disconnections gauge
total_number_of_disconnections{nodename="Linode"} 25
# HELP total_number_of_requests total number of requests
# TYPE total_number_of_requests gauge
total_number_of_requests{nodename="Linode"} 171
# HELP total_uptime_percentage total uptime percentage
# TYPE total_uptime_percentage gauge
total_uptime_percentage{nodename="Linode"} 99.60898158586997
# HELP total_uptime_seconds Total uptime of node in seconds
# TYPE total_uptime_seconds gauge
total_uptime_seconds{nodename="Linode"} 2.120476e+06
```