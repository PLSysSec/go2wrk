package structs

import (
	"github.com/kpister/go2wrk/logger"
	"net/http"
)

// TPSReport is a struct that contains the config file information.
type TPSReport struct {
	Routes         []Route
	Samples        int
	Connections    int
	MaxConnections int
	BreakMe        bool
	InitRoute      string
	MaxTestTime    float64
	Frequency      float64
	EndPercentage  float64
	Transport      *http.Transport
	Logger         *logger.Logger
}
