package structs

import "time"

type Response struct {
	Duration   int64
	Error      bool
	Size       int64
	Start      time.Time
	StatusCode int
}
