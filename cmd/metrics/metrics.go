package main

import (
	"context"
	"flag"
	"os"
	"runtime"

	metrics "meli-proxy/pkg/metrics"
	metricsserver "meli-proxy/pkg/server/metrics"
)

var redisAddrFlag := flag.String("redis", "", "Redis server")

func main() {

	flag.Parse()

	runtime.GOMAXPROCS(12)

	redisAddr := *redisAddrFlag

	if redisAddr == "" {
		redisAddr = os.Getenv("REDIS_SERVER")
		if redisAddr == "" {
			redisAddr = ":6379"
		}
	}

	ctx := context.Background()
	server := metricsserver.MetricsServer{
		Metrics: metrics.NewMetrics(ctx, redisAddr),
	}

	server.Run()

}
