package structs

import "sync"

// Bootstrap is a struct that stores the latencies of all the responses.
type Bootstrap struct {
	sync.Mutex
	List []int64
}
