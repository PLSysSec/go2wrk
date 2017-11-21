package structs

import "sync"

type Bootstrap struct {
	sync.Mutex
	MetricList []float64
}