package proxy

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"reverse-proxy/middleware"
	"reverse-proxy/pkg/metrics"
	"reverse-proxy/pkg/routes"
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

func (s *Server) processRequest(mapping string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %d - %s", s.port, r.RequestURI)

		ctx := r.Context()
		req, err := http.NewRequest(r.Method, mapping+r.RequestURI, nil)
		req = req.WithContext(ctx)
		req.Header.Add("Accept", "application/json")
		resp, err := s.Client.Do(req)

		if err != nil {
			log.Fatal("Errored when sending request to the server")
			log.Fatal(err)
			return
		}

		defer resp.Body.Close()

		for name, values := range resp.Header {
			w.Header()[name] = values
		}

		w.WriteHeader(resp.StatusCode)

		io.Copy(w, resp.Body)

	}
}

/*
Run allows to start Server
*/
func (s *Server) Run() {

	s.instance.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	testRoute := routes.RouteConfig{
		Path:      "/go/",
		Server:    "https://localhost",
		Limit:     200,
		Time:      time.Second,
		LimitType: routes.OriginIPUserAgent,
	}

	s.instance.HandleFunc("/go", middleware.TrackUser(testRoute, s.Metrics, middleware.IsAllowed(testRoute, s.Metrics,
		(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Go!"))
		}))))

	for _, route := range s.routes {
		log.Println("Adding route: ", route)
		s.instance.HandleFunc(route.Path, middleware.TrackUser(route, s.Metrics, middleware.IsAllowed(route, s.Metrics, s.processRequest(route.Server))))
	}

	listenPort := ":" + strconv.Itoa(s.port)

	log.Fatal(http.ListenAndServe(listenPort, middleware.WrapHandlerWithMetric(s.Metrics, s.instance)))

}

/*
NewProxy returns a new instance of the proxy Server
*/
func NewProxy(ctx context.Context, port int, routes []routes.RouteConfig, redisAddr string) (s Server) {

	s = Server{
		ctx:      ctx,
		port:     port,
		routes:   routes,
		Metrics:  metrics.NewMetrics(ctx, redisAddr),
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
