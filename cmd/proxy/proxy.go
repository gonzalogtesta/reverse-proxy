package main

import (
	"context"
	"flag"
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

	ctx := context.Background()
	// s := server.Proxy(ctx, *addr, routing)

	s := proxy.NewProxy(ctx, *addr, routing)
	s.Run()

}
