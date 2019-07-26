package main

import (
	"context"
	"flag"
	"time"

	"meli-proxy/cmd/proxy/server"
	"meli-proxy/pkg/routes"
)

var routing = []routes.RouteConfig{
	routes.RouteConfig{
		Path:   "/categories/",
		Server: "https://api.mercadolibre.com",
		Limit:  5000,
		Time:   time.Second,
	},
}

func main() {
	addr := flag.Int("addr", 8081, "HTTP network address")

	ctx := context.Background()
	s := server.Proxy(ctx, *addr, routing)
	s.Run()

}
