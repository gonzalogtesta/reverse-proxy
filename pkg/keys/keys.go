package keys

import (
	"fmt"
	"reverse-proxy/pkg/routes"
	"reverse-proxy/utils/ip"
	"net/http"
	"strings"
)

/*
Key constants
*/
const (
	StartKey                string = "user_request:"
	OriginIP                string = "IP:%s"
	DestinationPath         string = "Path:%s"
	OriginIPDestinationPath string = OriginIP + "_" + DestinationPath
	OriginIPUserAgent       string = OriginIP + "_UserAgent:%s"
)

/*
GenerateKey generates a key based on a http request
*/
func GenerateKey(route routes.RouteConfig, r *http.Request) string {

	key := ""
	switch route.LimitType {
	case routes.OriginIP:
		key = fmt.Sprintf(OriginIP, ip.GetIP(r))
	case routes.DestinationPath:
		key = fmt.Sprintf(DestinationPath, r.RequestURI)
	case routes.OriginIPDestinationPath:
		key = fmt.Sprintf(OriginIPDestinationPath, ip.GetIP(r), r.RequestURI)
	case routes.OriginIPUserAgent:
		key = fmt.Sprintf(OriginIPUserAgent, ip.GetIP(r), r.UserAgent())
	}

	return StartKey + key
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
