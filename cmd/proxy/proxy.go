package main

import (
	"context"
	"time"

	"meli-proxy/proxy"
)

var servingRoutes = map[string]string{
	"/categories/": "https://api.mercadolibre.com",
	"/cars/":       "http://apps-sysone-app.apps.us-east-2.online-starter.openshift.com",
}

var portToServer = map[int]map[string]string{
	8081: servingRoutes,
	8080: {"/": "http://apps-sysone-app.apps.us-east-2.online-starter.openshift.com"},
}

var config = proxy.RouteConfig{
	Path:   "/categories/",
	Server: "https://api.mercadolibre.com",
	Limit:  30,
	Time:   time.Second,
}

var routing = []proxy.RouteConfig{
	proxy.RouteConfig{
		Path:   "/categories/",
		Server: "https://api.mercadolibre.com",
		Limit:  30,
		Time:   time.Second,
	},
}

func main() {

	ctx := context.Background()
	s := proxy.Proxy(ctx, 8081, routing)
	s.Run()

}
