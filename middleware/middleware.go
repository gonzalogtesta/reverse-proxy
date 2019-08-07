package middleware

import (
	"reverse-proxy/pkg/keys"
	"reverse-proxy/pkg/metrics"
	"reverse-proxy/pkg/routes"
	"net/http"
	"time"
)

type Middleware struct {
}

/*
IsAllowed checks user limits configured in the RouteConfig
*/
func IsAllowed(route routes.RouteConfig, me metrics.Metrics, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, _ := me.GetForPeriod(keys.GenerateKey(route, r), route.Time); val >= route.Limit {
			http.Error(w, "Too Many Requests", 429)
			return
		}
		h(w, r)
	}
}

/*
TrackUser tracks an user request information
*/
func TrackUser(route routes.RouteConfig, me metrics.Metrics, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go me.Track(route, r)
		h(w, r)
	}
}

/*

 */
type metricResponseWriter struct {
	http.ResponseWriter
	statusCode int
	startTime  time.Time
}

func newMetricResponseWriter(w http.ResponseWriter) *metricResponseWriter {
	return &metricResponseWriter{w, http.StatusOK, time.Now()}
}

func (lrw *metricResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)

}

/*
WrapHandlerWithMetric allows to take metrics
*/
func WrapHandlerWithMetric(me metrics.Metrics, wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// go me.Track(r)
		go me.Hit(r)

		mrw := newMetricResponseWriter(w)
		wrappedHandler.ServeHTTP(mrw, r)

		statusCode := mrw.statusCode
		go me.SendCode(statusCode, mrw.startTime)
	})
}
