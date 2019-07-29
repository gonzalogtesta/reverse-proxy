package main

import (
	"context"
	"flag"
	"os"
	"runtime"

	metrics "meli-proxy/pkg/metrics"
	metricsserver "meli-proxy/pkg/server/metrics"
)

func main() {

	runtime.GOMAXPROCS(12)

	redisAddr := *flag.String("redis", "", "Redis server")

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
