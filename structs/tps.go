package structs

import (
	gofish "github.com/streddy/go-fish/structs"
	"net/http"
)

// TPSReport is a struct that contains the config file information.
type TPSReport struct {
	Routes        []gofish.Route
	Samples       int
	Connections   int
	InitRoute     string
	MinLatency    string
	MaxLatency    string
	MaxTestTime   float64
	Frequency     float64
	DropFreq      float64
	EndPercentage float64
	Transport     *http.Transport
	UseTransport  bool
}
