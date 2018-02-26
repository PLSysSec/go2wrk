package structs

import "time"

// Response is a struct that stores information about a http request response.
type Response struct {
	Duration   int64
	Error      bool
	Size       int64
	Start      time.Time
	StatusCode int
}
