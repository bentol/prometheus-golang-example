package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	curlErrorCollector = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "error_curl_total",
			Help: "Total curl request failed",
		},
		[]string{"vendor"},
	)
)

func init() {
	prometheus.MustRegister(curlErrorCollector)
}

func main() {
	http.Handle("/metrics", prometheus.Handler())

	// hit this api several time with query string vendor=something
	http.HandleFunc("/cek-resi", cekResi)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cekResi(w http.ResponseWriter, r *http.Request) {
	vendor := r.FormValue("vendor")

	// some logic here
	err := rand.Intn(2) == 0

	// handle error
	if err != false {
		// if error increment total error
		go RecordCurlError(vendor)
		w.Write([]byte("Failed to fetch"))
	} else {
		w.Write([]byte("Resi status: ok"))
	}

}

func RecordCurlError(vendor string) {
	curlErrorCollector.With(prometheus.Labels{"vendor": vendor}).Inc()
}
