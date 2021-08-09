package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var port = ":" + os.Getenv("PORT")

func main() {
	// initialize logger
	l := hclog.New(&hclog.LoggerOptions{
		Name: "Main",
	})

	// create a new server mux
	sm := mux.NewRouter()
	sm.Use(PrometheusMiddleware)

	sm.Path("/metrics").Handler(promhttp.Handler())

	// register handlers
	controller := NewController()
	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/status/{id}", controller.GetStatus)
	getRouter.HandleFunc("/generate", controller.GenerateOrder)

	// CORS
	goHandlers.CORS()

	// create http server
	s := &http.Server{
		Addr:         port,
		Handler:      sm,
		ErrorLog:     l.StandardLogger(&hclog.StandardLoggerOptions{}),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Info("Starting server...", "PORT", hclog.Fmt("%s", port))
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 10)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Info("Received terminate, graceful shutdown", "signal", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
