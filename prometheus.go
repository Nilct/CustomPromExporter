package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// Resets the gaugeVecs back to 0
// Ensures we start from a clean sheet
func (e *Exporter) resetGaugeVecs() {
	for _, m := range e.gaugeVecs {
		m.Reset()
	}
}

// Describe describes all the metrics ever exported by the Rancher exporter
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range e.gaugeVecs {
		m.Describe(ch)
	}
}

// Collect function, called on by Prometheus Client library
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	e.resetGaugeVecs() // Clean starting point

	var data, _ = e.gatherData(ch)

	e.processMetrics(data, ch)

	for _, m := range e.gaugeVecs {
		m.Collect(ch)
	}
	fmt.Printf("Collect\n")
}
