package routes

import (
	"testing"
	"time"
)

func TestReadRoute(t *testing.T) {

	t.Run("Parsing route", func(t *testing.T) {
		str := `
		{
			"path": "/path",
			"server": "http://test.com/",
			"limit": 1000,
			"time": "25m",
			"limitType": "OriginIP"
		  }
		`

		parsedRoute := parseRoute([]byte(str))

		route := RouteConfig{
			Path:      "/path",
			Server:    "http://test.com/",
			Limit:     1000,
			Time:      time.Minute * 25,
			LimitType: OriginIP,
		}

		if parsedRoute.Path != route.Path {
			t.Fatalf("Route path parsed is not the same, expecting: %s, parsed: %s", route.Path, parsedRoute.Path)
		}

		if parsedRoute.Server != route.Server {
			t.Fatalf("Route Server parsed is not the same, expecting: %s, parsed: %s", route.Server, parsedRoute.Server)
		}

		if parsedRoute.Limit != route.Limit {
			t.Fatalf("Route path parsed is not the same, expecting: %d, parsed: %d", route.Limit, parsedRoute.Limit)
		}

		if parsedRoute.Time != route.Time {
			t.Fatalf("Route path parsed is not the same, expecting: %d, parsed: %d", route.Time, parsedRoute.Time)
		}

		if parsedRoute.LimitType != route.LimitType {
			t.Fatalf("Route path parsed is not the same, expecting: %s, parsed: %s", route.LimitType, parsedRoute.LimitType)
		}
	})

}
