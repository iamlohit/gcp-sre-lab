package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Time spent serving HTTP requests.",
		Buckets: []float64{0.05, 0.1, 0.2, 0.5, 1, 2},
	}, []string{"method", "route", "status_code"})
)

type LoanData struct {
	LoanID string  `json:"loan_id"`
	Amount float64 `json:"amount"`
	Rate   float64 `json:"rate"`
	Status string  `json:"status"`
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	code := http.StatusOK

	// Simulate random delay/error for testing
	if rand.Float64() < 0.05 { // 5% error chance
		code = http.StatusInternalServerError
		http.Error(w, "Simulated error", code)
	} else {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond) // 0-200ms delay
		data := LoanData{"12345", 500000, 3.5, "approved"}
		json.NewEncoder(w).Encode(data)
	}

	duration := time.Since(start).Seconds()
	requestDuration.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", code)).Observe(duration)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/query", queryHandler)
	http.ListenAndServe(":8080", nil)
}
