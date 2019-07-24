package metrics

import (
	"context"
	"log"
	"net/http"

	"meli-proxy/metrics/consumer"
)

type server struct {
	db  string
	amq string
}

func (s *server) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

/*

 */
func Run() {

	ctx, cancel := context.WithCancel(context.Background())

	c := consumer.Consumer{}

	go c.Start(ctx)

	s := server{
		db:  "",
		amq: "",
	}
	http.HandleFunc("/", s.handler())
	log.Fatal(http.ListenAndServe(":1234", nil))

	defer cancel()
}
