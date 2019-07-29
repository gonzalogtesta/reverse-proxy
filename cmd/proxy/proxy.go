package main

import (
	"context"
	"flag"
	"os"
	"time"

	"meli-proxy/pkg/routes"
	"meli-proxy/pkg/server/proxy"
)

var routing = []routes.RouteConfig{
	routes.RouteConfig{
		Path:      "/categories/",
		Server:    "https://api.mercadolibre.com",
		Limit:     5000,
		Time:      time.Second,
		LimitType: routes.OriginIPDestinationPath,
	},
}

func main() {
	addr := flag.Int("addr", 8081, "HTTP network address")

	redisAddr := *flag.String("redis", "", "Redis server")

	if redisAddr == "" {
		redisAddr = os.Getenv("REDIS_SERVER")
		if redisAddr == "" {
			redisAddr = ":6379"
		}
	}

	ctx := context.Background()
	// s := server.Proxy(ctx, *addr, routing)

	s := proxy.NewProxy(ctx, *addr, routing, redisAddr)
	s.Run()

}
