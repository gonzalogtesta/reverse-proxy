package main

import (
	"context"
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

func main() {

	ctx := context.Background()
	/*
		// trap Ctrl+C and call cancel on the context
		ctx, cancel := context.WithCancel(ctx)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		defer func() {
			signal.Stop(c)
			cancel()
		}()
		go func() {
			select {
			case <-c:
				cancel()
			case <-ctx.Done():
			}
		}()
	*/
	s := proxy.Proxy(ctx, 8081, servingRoutes)
	s.Run()

}
