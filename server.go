package main

import (
	"context"
	"net/http"
	"sync"

	"github.com/miton18/go-warp10/instrumentation"
	log "github.com/sirupsen/logrus"
)

var server http.Server

func startServer(listen string, data chan instrumentation.Metrics) {
	m := http.NewServeMux()
	m.HandleFunc("/metrics", handler(data))

	server = http.Server{
		Addr:    listen,
		Handler: m,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.WithError(err).Fatal("cannot open HTTP server")
		}
	}()
}

func stopServer() error {
	return server.Shutdown(context.Background())
}

func handler(data chan instrumentation.Metrics) func(http.ResponseWriter, *http.Request) {
	metrics := instrumentation.Metrics{}
	lock := sync.RWMutex{}

	go func() {
		for newMetrics := range data {
			lock.Lock()
			metrics = newMetrics
			lock.Unlock()
		}

	}()

	return func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("","")
		w.WriteHeader(http.StatusOK)

		lock.RLock()
		_, err := w.Write([]byte(metrics.Get().Sensision()))
		lock.RUnlock()
		if err != nil {
			log.WithError(err).Error("cannot answer to client")
		}
	}
}
