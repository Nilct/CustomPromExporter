package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

func addMetrics() map[string]*prometheus.GaugeVec {
	gaugeVecs := make(map[string]*prometheus.GaugeVec)

	// Stack Metrics
	// gaugeVecs["imageDuration"] = prometheus.NewGaugeVec(
	// 	prometheus.GaugeOpts{
	// 		Namespace: "codingame",
	// 		Subsystem: "my_computer",
	// 		Name:      "docker_image_duration",
	// 		Help:      "Docker image existence duration (in be decided)",
	// 	}, []string{"repository", "tag"})
	// gaugeVecs["imageSize"] = prometheus.NewGaugeVec(
	// 	prometheus.GaugeOpts{
	// 		Namespace: "codingame",
	// 		Subsystem: "my_computer",
	// 		Name:      "docker_image_size",
	// 		Help:      "Docker image size (in Mo)",
	// 	}, []string{"repository", "tag"})
	gaugeVecs["containerMemory"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "codingame",
			Subsystem: "my_computer",
			Name:      "docker_container_memory",
			Help:      "Container usage memory",
		}, []string{"image", "id"})
	return gaugeVecs
}

// setMetrics - Logic to set the state of a system as a gauge metric
func (e *Exporter) setMetrics(image, id string, memory uint64) error {
	//fmt.Printf("Update %s:%s with duration: %d size: %d\n", repository, tag, duration, size)
	//e.gaugeVecs["imageDuration"].With(prometheus.Labels{"repository": repository, "tag": tag}).Set(float64(duration))
	//e.gaugeVecs["imageSize"].With(prometheus.Labels{"repository": repository, "tag": tag}).Set(float64(size))
	e.gaugeVecs["containerMemory"].With(prometheus.Labels{"image": image, "id": id}).Set(float64(memory))
	return nil
}
