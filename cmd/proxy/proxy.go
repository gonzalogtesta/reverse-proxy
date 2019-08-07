package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"reverse-proxy/pkg/routes"
	"reverse-proxy/pkg/server/proxy"
)

var routing = []routes.RouteConfig{
	routes.RouteConfig{
		Path:      "/test/",
		Server:    "https://localhost",
		Limit:     5000,
		Time:      time.Second,
		LimitType: routes.OriginIPDestinationPath,
	},
}

var addrFlag = flag.Int("addr", 8081, "HTTP network address")
var redisAddrFlag = flag.String("redis", "", "Redis server")
var configFile = flag.String("config", "", "Config file")

func main() {
	flag.Parse()

	redisAddr := *redisAddrFlag

	if redisAddr == "" {
		redisAddr = os.Getenv("REDIS_SERVER")
		if redisAddr == "" {
			redisAddr = ":6379"
		}
	}

	if *configFile != "" {
		fmt.Println("Reading config file: ", *configFile)
		routing = append(routing, routes.ReadFileRoute(*configFile))
	}

	ctx := context.Background()

	s := proxy.NewProxy(ctx, *addrFlag, routing, redisAddr)
	s.Run()

}
