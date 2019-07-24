package proxy

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"meli-proxy/statistics"
)

/*
Server struct
*/
type Server struct {
	port     int
	routes   map[string]string
	metrics  statistics.Metrics
	instance *http.ServeMux
	client   *http.Client
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
	resp, err := s.client.Do(req)

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
		go s.metrics.SendMetrics(r)

		s.getRequest(w, r, mapping)
	}
}

func (s *Server) isAllowed(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, _ := s.metrics.GetForPeriod(r, time.Second*30); val >= 300 {
			http.Error(w, "Too Many Requests", 429)
			// 429 Too Many Requests
			return
		}
		h(w, r)
	}
}

func (s *Server) metricRequest(h http.HandlerFunc) http.HandlerFunc {
	fmt.Println("Middleware ON")
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		fmt.Println("Sending Metrics")
		go s.metrics.Track(r)

	}
}

/*
Run allows to start Server
*/
func (s *Server) Run() {

	s.instance.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

	})

	s.instance.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		data, _ := s.metrics.Get(time.Second * 30)
		fmt.Println("data")
		fmt.Println(data)
		buf, _ := json.Marshal(data)
		w.Write(buf)
		// fmt.Println(string(data))
		/*
			for _, val := range data {
				fmt.Println(val)
				w.Write([]byte(val))
			}
		*/
	})

	s.instance.HandleFunc("/go", s.isAllowed(s.metricRequest(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go!"))
	})))

	for route, mapping := range s.routes {
		s.instance.HandleFunc(route, s.processRequest(mapping))
	}

	listenPort := ":" + strconv.Itoa(s.port)

	if err := http.ListenAndServe(listenPort, s.instance); err != nil {
		panic(err) // don't panic
	}
}

/*
Proxy allows to generate an instance of a proxy server for an specified port.

*/
func Proxy(ctx context.Context, port int, routes map[string]string) Server {

	s := Server{
		ctx:      ctx,
		port:     port,
		routes:   routes,
		metrics:  (&statistics.Metrics{}).Connect(ctx),
		instance: http.NewServeMux(),
		client: &http.Client{
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
