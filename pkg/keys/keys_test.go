package keys

import (
	"net/http"
	"testing"
	"time"

	"meli-proxy/pkg/routes"
)

func TestGenerateKeyUsingRouteOriginIP(t *testing.T) {
	t.Run("GenerateKey Using Route OriginIP", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/test", nil)

		route := routes.RouteConfig{
			Path:      "/test",
			Server:    "http://test_server:8080/",
			Limit:     200,
			Time:      time.Minute,
			LimitType: routes.OriginIP,
		}

		request.RemoteAddr = "8.8.8.8"

		key := GenerateKey(route, request)
		expected := "user_request:IP:8.8.8.8"
		if key != expected {
			t.Fatalf("Invalid key expecting: %s, generated: %s", expected, key)
		}
	})

	t.Run("GenerateKey Using Route OriginIP with X-Forwarded-For", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/test", nil)

		route := routes.RouteConfig{
			Path:      "/test",
			Server:    "http://test_server:8080/",
			Limit:     200,
			Time:      time.Minute,
			LimitType: routes.OriginIP,
		}

		request.RemoteAddr = ""
		request.Header.Set("X-Forwarded-For", "9.9.9.9")

		key := GenerateKey(route, request)
		expected := "user_request:IP:9.9.9.9"
		if key != expected {
			t.Fatalf("Invalid key expecting: %s, generated: %s", expected, key)
		}
	})

	t.Run("GenerateKey Using Route DestinationPath", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/test", nil)

		route := routes.RouteConfig{
			Path:      "/test",
			Server:    "http://test_server:8080/",
			Limit:     200,
			Time:      time.Minute,
			LimitType: routes.DestinationPath,
		}

		request.RemoteAddr = "8.8.8.8"
		request.RequestURI = "/test"

		key := GenerateKey(route, request)
		expected := "user_request:Path:/test"
		if key != expected {
			t.Fatalf("Invalid key expecting: %s, generated: %s", expected, key)
		}
	})

	t.Run("GenerateKey Using Route IP and DestinationPath", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/test", nil)

		route := routes.RouteConfig{
			Path:      "/test",
			Server:    "http://test_server:8080/",
			Limit:     200,
			Time:      time.Minute,
			LimitType: routes.OriginIPDestinationPath,
		}

		request.RemoteAddr = "10.10.10.10"
		request.RequestURI = "/test"

		key := GenerateKey(route, request)
		expected := "user_request:IP:10.10.10.10_Path:/test"
		if key != expected {
			t.Fatalf("Invalid key expecting: %s, generated: %s", expected, key)
		}
	})

	t.Run("GenerateKey Using RouteIP and UserAgent", func(t *testing.T) {

		request, _ := http.NewRequest("GET", "/test", nil)

		route := routes.RouteConfig{
			Path:      "/test",
			Server:    "http://test_server:8080/",
			Limit:     200,
			Time:      time.Minute,
			LimitType: routes.OriginIPUserAgent,
		}

		request.RemoteAddr = "11.11.11.11"
		request.RequestURI = "/test"
		request.Header.Set("User-Agent", "golang")

		key := GenerateKey(route, request)
		expected := "user_request:IP:11.11.11.11_UserAgent:golang"
		if key != expected {
			t.Fatalf("Invalid key expecting: %s, generated: %s", expected, key)
		}
	})
}
