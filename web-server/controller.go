package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/xid"
)

type Controller struct {
	l  hclog.Logger
	rc *RabbitMQClient
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

var totalRequest = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_total",
		Help: "Number of requests",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "status of HTTP response",
	},
	[]string{"status"},
)

func NewController() *Controller {
	rc := NewClient()
	log := hclog.New(&hclog.LoggerOptions{
		Name: "Handler",
	})

	return &Controller{l: log, rc: rc}
}

func (c *Controller) GenerateOrder(rw http.ResponseWriter, r *http.Request) {
	c.l.Info("generete order")
	id := xid.New().String()
	c.rc.SendMessage([]byte(id))
	fmt.Fprintln(rw, "Welcome!", id)
}

func (c *Controller) GetStatus(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	c.l.Info("get status", "ID", id)

	fmt.Fprintln(rw, "Welcome! ", id)

}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) writeHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequest.WithLabelValues(path).Inc()
	})
}

func init() {
	prometheus.Register(totalRequest)
	prometheus.Register(responseStatus)
}
