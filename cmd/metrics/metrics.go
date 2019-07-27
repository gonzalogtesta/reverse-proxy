package main

import (
	"context"

	metrics "meli-proxy/pkg/metrics"
	metricsserver "meli-proxy/pkg/server/metrics"
)

func main() {

	ctx := context.Background()
	server := metricsserver.MetricsServer{
		Metrics: metrics.NewMetrics(ctx),
	}

	server.Run()

}
