package keys

import (
	"meli-proxy/pkg/routes"
	"meli-proxy/utils/ip"
	"net/http"
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
