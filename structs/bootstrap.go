package structs

import "sync"

type Bootstrap struct {
	sync.Mutex
	List []int64
}
