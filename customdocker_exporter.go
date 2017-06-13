package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "mydockerexporter" // Used to prepand Prometheus metrics created by this exporter.
)

var (
	metricsPath   = "/metrics"           // Path under which to expose metrics
	listenAddress = "192.168.1.129:8083" // Address on which to expose metrics
)

func main() {
	flag.Parse()

	// Register internal metrics used for tracking the exporter performance
	//measure.Init()
	InitScraping()
	// Register a new Exporter
	Exporter := newExporter()

	// Register Metrics from each of the endpoints
	// This invokes the Collect method through the prometheus client libraries.
	prometheus.MustRegister(Exporter)

	// Setup HTTP handler
	http.Handle(metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		                <head><title>My custom docker exporter</title></head>
		                <body>
		                   <h1>docker exporter</h1>
		                   <p><a href='` + metricsPath + `'>Metrics</a></p>
		                   </body>
		                </html>
		              `))
	})
	log.Printf("Starting Server on port %s and path %s", listenAddress, metricsPath)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
