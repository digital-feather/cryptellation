package health

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

const HealthDefaultPort = 9000

type Health struct {
	isReady atomic.Value
}

func New() *Health {
	var h Health

	h.isReady.Store(false)

	return &h
}

func (h *Health) Ready(isReady bool) {
	h.isReady.Store(isReady)
}

func (h *Health) Serve() {
	port, err := strconv.Atoi(os.Getenv("CRYPTELLATION_HEALTH_PORT"))
	if err != nil {
		log.Println("error when parsing health port:", err)
		log.Println("setting to default:", HealthDefaultPort)
		port = HealthDefaultPort
	}

	http.HandleFunc("/liveness", h.liveness())
	http.HandleFunc("/readiness", h.readiness())

	url := fmt.Sprintf(":%d", port)
	log.Println("Starting: HTTP Health Listener")
	log.Fatal(http.ListenAndServe(url, nil))
}

func (h *Health) liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Health) readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if !h.isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
