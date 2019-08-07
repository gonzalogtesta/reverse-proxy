package metricsserver

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"html/template"
	"net/http"
	"time"

	"reverse-proxy/pkg/keys"
	"reverse-proxy/pkg/metrics"
)

type MetricsServer struct {
	Metrics metrics.Metrics
}

type DataPercentile struct {
	Key  string
	Data [][]float64
}

type DataCount struct {
	Key  string
	Data [][]int64
}

type Data struct {
	Hits         [][]int64
	Response200  [][]int64
	Response404  [][]int64
	Response429  [][]int64
	Response500  [][]int64
	Percentile90 [][]float64
	Percentile95 [][]float64
	Percentile99 [][]float64
	Percentiles  map[string][][]float64
	Counters     map[string][][]int64
}

func (s *MetricsServer) metricsRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricKeys := r.URL.Query()["metrics"]
		fmt.Println("Metrics: ", metricKeys)

		timeFrame := r.URL.Query().Get("time")

		generic, simple := keys.GroupKeys(metricKeys)

		items := s.Metrics.GetKeys(generic)
		simple = append(simple, items...)

		data := make(map[string]interface{})
		dur, err := time.ParseDuration(timeFrame)
		if err != nil {
			dur = time.Second * 30
			fmt.Println("Unable to parse duration")
			fmt.Println(err)
		}
		for _, key := range simple {
			data[key], _ = s.Metrics.GetSerie(key, dur)
		}

		buf, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf)
	}
}

func (s *MetricsServer) metricsPercentilesRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		percentileStr := r.URL.Query().Get("percentile")
		fmt.Println("Percentile: ", percentileStr)

		percentile, err := strconv.Atoi(percentileStr)
		if err != nil {
			http.Error(w, "Invalid percentile", http.StatusNotAcceptable)
			return
		}

		timeFrame := r.URL.Query().Get("time")
		duration, err := time.ParseDuration(timeFrame)
		if err != nil {
			duration = time.Second * 30
		}
		metric := r.URL.Query().Get("metric")
		data, _ := s.Metrics.GetPercentile(metric, percentile, duration)

		buf, _ := json.Marshal(data)
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf)
	}
}

func (s *MetricsServer) metricsHTML(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl = template.Must(template.ParseFiles("statics/layout.html"))
		tmpl.Execute(w, Data{})
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
