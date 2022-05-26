# uptime-robot-exporter
Exports [Uptime Robot](https://uptimerobot.com/) monitor statuses to [Prometheus](https://prometheus.io/)

## Running

Set the following environment variables:
- INTERVAL_POLLING: the interval time to poll Uptime Robot API in `time.Duration` format (10s, 30s, 1m, etc). Keep in mind the API rate limits (FREE plan : 10 req/min, PRO plan : monitor limit * 2 req/min with maximum value 5000 req/min)
- UPTIME_ROBOT_API_KEY: this is the API key for Uptime Robot (supports both per monitor and for all monitors)

Build & Run:
```bash
$: cd src/
$: go get
$: go build
$: ./uptime-robot-exporter
```

## Metrics

Currently, only status metric is exported as `uptime_robot_status` with the tags "name", "url", "type", "sub_type", "port".
According to [Uptime Robot API docs](https://uptimerobot.com/api/), the value of the status is mapped as follows:	
- 0 - paused
- 1 - not checked yet
- 2 - up
- 8 - seems down
- 9 - down

## Development

`docker-compose` can be used to create a development stack once .env file is created with appropiate UPTIME_ROBOT_API_KEY and POLLING_INTERVAL variables:

```bash
cd development/
docker-compose up # -d if you want to run in dettached mode
```

Then, prometheus should be reachable at http://localhost:9090/ and configured with uptime-robot-exporter as target.

## Libraries

- [zerolog](https://github.com/rs/zerolog) for logging
- JSON to struct extracted with [json-to-go](https://mholt.github.io/json-to-go/)

## TO-DO

- [ ] Control API rate limits and other status codes
- [ ] Add Tests
- [ ] More metrics based in Uptime Robot API response
- [ ] Helm Chart for k8s deployment
- [ ] Improve logging
- [ ] Add Uptime Robot log events
