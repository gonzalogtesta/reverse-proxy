package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"meli-proxy/middleware"
	"meli-proxy/pkg/metrics"
	"meli-proxy/pkg/routes"
	"meli-proxy/utils/ip"
)

/*
Server struct
*/
type Server struct {
	port     int
	routes   []routes.RouteConfig
	Metrics  metrics.Metrics
	instance *http.ServeMux
	Client   *http.Client
	ctx      context.Context
}

func (s *Server) getRequest(w http.ResponseWriter, r *http.Request, mapping string) {
	fmt.Println("Request:")
	fmt.Println(s.port)

	ctx, cancel := context.WithCancel(context.TODO())
	timer := time.AfterFunc(5*time.Second, func() {
		cancel()
	})

	req, err := http.NewRequest(r.Method, mapping+r.RequestURI, nil)
	req = req.WithContext(ctx)
	req.Header.Add("Accept", "application/json")
	resp, err := s.Client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		fmt.Println(err)
		return
	}
	timer.Stop()

	defer resp.Body.Close()

	for name, values := range resp.Header {
		w.Header()[name] = values
	}

	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
	resp.Body.Close()

}

func (s *Server) processRequest(mapping string) http.HandlerFunc {
	// c := make(chan http.ResponseWriter)
	return func(w http.ResponseWriter, r *http.Request) {
		// go s.metrics.SendMetrics(r)

		s.getRequest(w, r, mapping)
	}
}

func (s *Server) isAllowed(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, _ := s.Metrics.GetForPeriod("user_request:"+ip.GetIP(r), time.Second*30); val >= 5000 {
			http.Error(w, "Too Many Requests", 429)
			// 429 Too Many Requests
			return
		}
		h(w, r)
	}
}

/*
func (s *Server) metricRequest(h http.HandlerFunc) http.HandlerFunc {
	fmt.Println("Middleware ON")
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		fmt.Println("Sending Metrics")
		go s.Metrics.Track(r)
		go s.Metrics.Hit(r)

	}
}
*/
type Data struct {
	Hits        [][]int64
	Response200 [][]int64
	Response404 [][]int64
	Response429 [][]int64
	Response500 [][]int64
}

/*
Run allows to start Server
*/
func (s *Server) Run() {

	s.instance.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	s.instance.HandleFunc("/go", s.isAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go!"))
	}))

	for _, route := range s.routes {
		s.instance.HandleFunc(route.Path, middleware.IsAllowed(route, s.Metrics, s.processRequest(route.Server)))
	}

	listenPort := ":" + strconv.Itoa(s.port)

	if err := http.ListenAndServe(listenPort, middleware.WrapHandlerWithMetric(s.Metrics, s.instance)); err != nil {
		panic(err) // don't panic
	}
}

/*
Proxy allows to generate an instance of a proxy server for an specified port.

*/
func Proxy(ctx context.Context, port int, routes []routes.RouteConfig) Server {

	s := Server{
		ctx:      ctx,
		port:     port,
		routes:   routes,
		Metrics:  metrics.NewMetrics(ctx),
		instance: http.NewServeMux(),
		Client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	return s
}
