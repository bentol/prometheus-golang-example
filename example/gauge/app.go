package main

import (
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	goroutinesGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "my_app_goroutines_count",
			Help: "number of goroutines that currently exist",
		},
		[]string{"hostname"},
	)
)

func init() {
	prometheus.MustRegister(goroutinesGauge)
}

func observer() {
	hostname, _ := os.Hostname()

	for {
		value := runtime.NumGoroutine()
		goroutinesGauge.With(prometheus.Labels{"hostname": hostname}).Set(float64(value))
		time.Sleep(1 * time.Second)
	}
}

func main() {
	go observer()
	http.Handle("/metrics", prometheus.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
