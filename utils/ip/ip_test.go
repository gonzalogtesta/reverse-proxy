package ip

import (
	"net/http"
	"testing"
)

func TestIP(t *testing.T) {

	t.Run("Parse ipv4 from Remote Addr", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/test", nil)
		request.RemoteAddr = "8.8.8.8"
		ip := GetIP(request)
		expected := "8.8.8.8"
		if ip != expected {
			t.Fatalf("IP expected was not the same as the request, expected: %s, actual: %s", expected, ip)
		}
	})

	t.Run("Parse local ipv6 from Remote Addr", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/test", nil)
		request.RemoteAddr = "[::1]"
		ip := GetIP(request)
		expected := "[::1]"
		if ip != expected {
			t.Fatalf("IP expected was not the same as the request, expected: %s, actual: %s", expected, ip)
		}
	})

	t.Run("Parse ipv6 from Remote Addr", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/test", nil)
		request.RemoteAddr = "[1200:0000:AB00:1234:0000:2552:7777:1313]"
		ip := GetIP(request)
		expected := "[1200:0000:AB00:1234:0000:2552:7777:1313]"
		if ip != expected {
			t.Fatalf("IP expected was not the same as the request, expected: %s, actual: %s", expected, ip)
		}
	})

}
