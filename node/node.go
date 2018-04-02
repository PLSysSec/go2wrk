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
// add threshold and tails to the params
func Run(tps structs.TPSReport, outputDirectory string, outputIteration int) {
	var channels []chan *structs.Response
	for i := 0; i < len(tps.Routes); i++ {
		// we might not need this anymore
		channels = append(channels, make(chan *structs.Response, int64(tps.MaxTestTime*tps.Frequency)*int64(tps.Connections)))
	}

	// we might not need this anymore
	// shared response metric collector and corresponding lock
	metrics := structs.Bootstrap{
		List:          make([]int64, 0),
		Converged:     false,
		Samples:       tps.Samples,
		EndPercentage: tps.EndPercentage,
	}
	waitGroup := &sync.WaitGroup{}
	start := time.Now()

	// connection.Init(tps)

	for i := 0; i < tps.Connections; i++ {
		// add threshold and tails to the params
		go connection.Start(tps, channels, start, &metrics, waitGroup)
		waitGroup.Add(1)
	}
	// doing this in main
	go (&metrics).Start() // start bootstrapping

	waitGroup.Wait()
	fmt.Println()

	for i, route := range tps.Routes {
		close(channels[i])
		// update export to deal with tails
		stats.Export(channels[i], i, outputIteration, route.Url, outputDirectory)
	}
	fmt.Printf("Response numbers: %d\n", len(metrics.List))
}
