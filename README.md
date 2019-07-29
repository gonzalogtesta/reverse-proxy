# meli-proxy

Implementation of a Reverse Proxy in Go, using Redis Time Series for tracking metrics.

It's divided in two applications, the proxy server and a metrics server.

## Instalation

### Docker

The application includes a docker-compose file, so it's possible to run all the applications with one command, prior to that the servers images must be generated:

For building the proxy server:
> docker build --no-cache -t proxy-server -f Dockerfile.proxy .


For building the metrics server:
> docker build --no-cache -t metrics-server -f Dockerfile.metrics . 


Once builded both images, the following command can be executed:

> docker-compose up --build

Also, it's possible to run both apps using golang directly, in order to do that, it's required to startup a Redis instance with RedisTimeSeries module.
Or using the image: *redislabs/redistimeseries*.

Install dependencies:

> go get -d -v ./...

Proxy:

> go run ./cmd/proxy/proxy.go

Metrics:

> go run ./cmd/metrics/metrics.go

For both applications the *redis* parameter can be provided to specify the host and port of Redis server.

## Usage Proxy

The proxy server allows to make request to different APIs, these must be configurated prior to proxy startup.

Application can be started in any port, always specifying the parameter *addr*. Eg:

> go run ./cmd/proxy/proxy.go --addr 8081

It allows to limit request rate based on:
- Origin IP
- Destination Path
- Origin IP + Destination Path
- Origin IP + User Agent

These limits can be configured by type and a # of requests per time, these both variables also can be configurable.

When limits are applied, if the limit is reached a 429 status code is returned. 
Rate limiter use an algorithm that calculates if the user exceeded the limit in a time window, taking in count the number of request performed in that time.

### Proxy request examples:

> curl http://localhost:8081/categories/MLA1000

> curl http://localhost:8081/cars


## Usage Metrics

Metrics server provides 2 endpoints for retrieving metrics:
 - /metrics
 - /metrics/percentiles
 
 After starting up Metrics server, it will start listening to port 4000. 
 
### /metrics endpoint

Retrieves the metrics stored.

Path: /metrics
Parameters:
  - metrics : list of metrics to be retrieved, it allows to specify a wildcard * to retrieve metrics without using the full name 
  - time : time specifies the window of time to retrieve information, eg: 1h, 30m, 60s...
  
  
Returns: a map with the metric as key and an array containing the serie. Eg:

```json
{
  "response_200": [
    [
      1564436046000,
      1
    ],
    [
      1564440156000,
      1
    ],
    [
      1564440157000,
      1
    ]
  ]
}
```
  
### /metrics/percentiles endpoint

Allows to calculate percentiles for a metric

Path: /metrics/percentiles
Parameters:
  - percentile : integer value of the percentile to be retrieved, eg: 80, 90, 95, 99
  - metric : metric name to calculate the percentile, eg.: response_200
  - time : time specifies the window of time to retrieve information, eg: 1h, 30m, 60s...
  
Returns: a list of the items calculated for the percentil

```json
[
  [
    1564440157000,
    0.30218040000000007
  ],
  [
    1564440156000,
    0.44138220000000006
  ]
]
```

## Examples

- Retrieve all hits in the last 24 hours: http://localhost:4000/metrics?metrics=hits&time=24h

- Retrieve 200 status responses in the last 24 hours: http://localhost:4000/metrics?metrics=response_200&time=24h

- Retrieve metrics for 200 and 404 status codes in the last 24 hours: http://localhost:4000/metrics?metrics=response_404&metrics=response_200&time=24h

- Retrieve metrics for all the *response_* using wildcard *, for the last 24 hours: http://localhost:4000/metrics?metrics=response_*&time=24h

- Retrieve 90th percentile for metric response_200 in the last 24 hs: http://localhost:4000/metrics/percentiles?time=24h&percentile=90&metric=response_200


## Live demo

There is an instance of both applications running in an AWS:

 - Proxy server: http://ec2-3-17-190-52.us-east-2.compute.amazonaws.com:8080/categories/MLA1000
 - Metrics server: http://ec2-3-17-190-52.us-east-2.compute.amazonaws.com:4000/metrics/percentiles?percentile=90&time=1h&metric=response_200


## Dashboard

Metrics server includes an example dashboard using Highcharts, that include:

 - Hits
 - Response 200
 - Response 404
 - Response 429
 - Response 500
 - Percentile 90th
 - Percentile 95th
 - Percentile 99th
 
Dashboard URL: http://ec2-3-17-190-52.us-east-2.compute.amazonaws.com:4000/metrics/html?time=24h
