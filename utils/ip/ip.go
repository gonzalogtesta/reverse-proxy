package ip

import (
	"net/http"
	"strings"
)

/*
GetIP gets an ip address from a http.Request
*/
func GetIP(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Forwarded-For")
	if ip == "" {
		if strings.Count(r.RemoteAddr, ":") < 2 {
			ip = strings.Split(r.RemoteAddr, ":")[0]
		} else {
			ip = strings.Split(r.RemoteAddr, "]")[0] + "]"
		}
	}
	return
}
