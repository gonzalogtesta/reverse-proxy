package metricsserver

import (
	"encoding/json"
	"fmt"
	"log"

	"html/template"
	"net/http"
	"time"

	"meli-proxy/pkg/keys"
	"meli-proxy/pkg/metrics"
)

type MetricsServer struct {
	Metrics metrics.Metrics
}

type Data struct {
	Hits        [][]int64
	Response200 [][]int64
	Response404 [][]int64
	Response429 [][]int64
	Response500 [][]int64
}

func (s *MetricsServer) metricsRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricKeys := r.URL.Query()["metrics"]
		fmt.Println("Metrics: ", metricKeys)

		timeFrame := r.URL.Query().Get("time")

		generic, simple := keys.GroupKeys(metricKeys)
		fmt.Println("Generic: ", generic)
		fmt.Println("Simple: ", simple)
		items := s.Metrics.GetKeys(generic)
		fmt.Println("Items: ", items)
		simple = append(simple, items...)

		data := make(map[string]interface{})
		dur, err := time.ParseDuration(timeFrame)
		if err != nil {
			dur = time.Second * 30
		}
		for _, key := range simple {
			data[key], _ = s.Metrics.GetSerie(key, dur)
		}
		//data, _ := s.Metrics.Get(time.Second * 30)
		//s.Metrics.GetSerie("hits", time.Second*30)

		buf, _ := json.Marshal(data)
		w.Write(buf)
	}
}

func (s *MetricsServer) metricsPercentilesRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		data := s.Metrics.GetPercentile(90, time.Minute*30)

		buf, _ := json.Marshal(data)
		w.Write(buf)
	}
}

func (s *MetricsServer) metricsHTML(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, _ := s.Metrics.GetSerie("hits", time.Second*30)
		resp200, _ := s.Metrics.GetSerie("response_200", time.Second*30)
		resp404, _ := s.Metrics.GetSerie("response_404", time.Second*30)
		resp429, _ := s.Metrics.GetSerie("response_429", time.Second*30)
		resp500, _ := s.Metrics.GetSerie("response_500", time.Second*30)
		d := Data{
			Hits:        data,
			Response200: resp200,
			Response404: resp404,
			Response429: resp429,
			Response500: resp500,
		}
		tmpl = template.Must(template.ParseFiles("statics/layout.html"))
		tmpl.Execute(w, d)
	}
}

func (s *MetricsServer) Run() {
	tmpl := template.Must(template.ParseFiles("statics/layout.html"))

	// http.HandleFunc("/", s.home)
	http.Handle("/metrics/percentiles", s.metricsPercentilesRoute())
	http.Handle("/metrics", s.metricsRoute())
	http.Handle("/metrics/html", s.metricsHTML(tmpl))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)
}
