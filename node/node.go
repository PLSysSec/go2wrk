package node

import (
	"github.com/kpister/go2wrk/connection"
	"github.com/kpister/go2wrk/stats"
	"github.com/kpister/go2wrk/structs"

	"fmt"
	"sync"
	"time"
)

// Warmup performs a short warmup on the server. These response are not recorded.
func Warmup(tps structs.TPSReport) {
	waitGroup := &sync.WaitGroup{}
	start := time.Now()
	for i := 0; i < tps.Connections; i++ {
		go connection.Warmup(tps, start, waitGroup)
		waitGroup.Add(1)
	}
	waitGroup.Wait()
	fmt.Println()
}

// Run will create connections that fire requests at the server. Then it creates the output.
func Run(tps structs.TPSReport, outputDirectory string) {
	var channels []chan *structs.Response
	for i := 0; i < len(tps.Routes); i++ {
		// TODO make this number meaningful
		channels = append(channels, make(chan *structs.Response, int(tps.TestTime)*tps.Connections*10))
	}

	// shared response metric collector and corresponding lock
	metrics := structs.Bootstrap{
		List: make([]int64, 0),
	}
	waitGroup := &sync.WaitGroup{}
	start := time.Now()

	for i := 0; i < tps.Connections; i++ {
		go connection.Start(tps, channels, start, &metrics, waitGroup)
		waitGroup.Add(1)
	}

	waitGroup.Wait()
	fmt.Println()

	for i, route := range tps.Routes {
		close(channels[i])
		stats.Export(channels[i], i, route.Url, outputDirectory)
	}
	fmt.Printf("Response numbers: %d\n", len(metrics.List))
}
