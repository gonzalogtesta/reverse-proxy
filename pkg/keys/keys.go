package keys

import (
	"meli-proxy/pkg/routes"
	"meli-proxy/utils/ip"
	"net/http"
	"strings"
)

func remoteAddr(r *http.Request) string {
	return "OK"
}

/*

const (
	REMOTEADDR        func = remoteAddr
	REQUESTURI string = "RequestURI"
)

*/

/*
GenerateKey generates a key based on a http request
*/
func GenerateKey(r *http.Request, route routes.RouteConfig) string {

	return "user_request:" + ip.GetIP(r) //fmt.Sprintf("a %s", "string")
}

/*
GroupKeys group the keys parsed in groups of generic or simple keys
*/
func GroupKeys(keys []string) (generic, simple []string) {
	for _, val := range keys {
		if strings.Contains(val, "*") {
			generic = append(generic, val)
		} else {
			simple = append(simple, val)
		}
	}

	return generic, simple
}
