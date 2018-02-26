package structs

import "net/http"

// TPSReport is a struct that contains the config file information.
type TPSReport struct {
	Routes      []Route
	Samples     int
	Connections int
	TestTime    float64
	Frequency   float64
	Latency     float64
	Transport   *http.Transport
}
