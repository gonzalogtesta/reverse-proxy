package middleware

import (
	"meli-proxy/pkg/keys"
	"meli-proxy/pkg/metrics"
	"meli-proxy/pkg/routes"
	"net/http"
	"time"
)

type Middleware struct {
}

/*
func (m *Middleware) oldprocessRequest(mapping string, me *metrics.Metrics, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, _ := me.GetForPeriod(r, time.Second*30); val >= 300 {
			http.Error(w, "Too Many Requests", 429)
			// 429 Too Many Requests
			return
		}
		h(w, r)
	}
}
*/
func IsAllowed(route routes.RouteConfig, me metrics.Metrics, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if val, _ := me.GetForPeriod(keys.GenerateKey(r, route), time.Second*route.Time); val >= route.Limit {
			http.Error(w, "Too Many Requests", 429)
			// 429 Too Many Requests
			return
		}
		h(w, r)
	}
}

/*
func (m *Middleware) metricRequest(me *metrics.Metrics, h http.HandlerFunc) http.HandlerFunc {
	fmt.Println("Middleware ON")
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		fmt.Println("Sending Metrics")
		go me.Track(r)
		go me.Hit(r)

	}
}
*/
type metricResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewMetricResponseWriter(w http.ResponseWriter) *metricResponseWriter {
	return &metricResponseWriter{w, http.StatusOK}
}

func (lrw *metricResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)

}
func WrapHandlerWithMetric(me metrics.Metrics, wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		go me.Track(r)
		go me.Hit(r)

		mrw := NewMetricResponseWriter(w)
		wrappedHandler.ServeHTTP(mrw, r)

		statusCode := mrw.statusCode
		go me.SendCode(statusCode)
	})
}
