package node

import (
	"github.com/kpister/go2wrk/connection"
	"github.com/kpister/go2wrk/stats"
	"github.com/kpister/go2wrk/structs"

	"time"
	"net/http"
	"fmt"
	"sync"
)

// Warmup performs a short warmup on the server. 
func Warmup(tps structs.TPSReport, index int) int64{
	waitGroup := &sync.WaitGroup{}
	warmupData := make(chan *structs.Response, 100 * int64(tps.Connections))

	tr := &http.Transport{
		MaxIdleConns:       0,
		IdleConnTimeout:    10 * time.Second,
	}
	client := &http.Client{Transport: tr}
	for i := 0; i < tps.Connections; i++ {
		go connection.Start(client, tps.Routes[index], tps.Frequency, warmupData, nil, waitGroup)
		waitGroup.Add(1)
	}
	waitGroup.Wait()
	close(warmupData)
	return stats.FindThreshold(warmupData)
}

// Barrage will create connections that fire requests at the server. Then it creates the output.
func Barrage(tps structs.TPSReport, outputDirectory string, outputIteration int) {
	var channels []chan *structs.Response
	for i := 0; i < len(tps.Routes); i++ {
		channels = append(channels, make(chan *structs.Response, int64(1000*tps.Frequency)*int64(tps.Connections)))
	}
	tr := &http.Transport{
		MaxIdleConns:       0,
		IdleConnTimeout:    10 * time.Second,
	}
	client := &http.Client{Transport: tr}

	// shared response metric collector and corresponding lock
	metrics := structs.Bootstrap{
		List:          make([]int64, 0),
		Converged:     false,
		Samples:       tps.Samples,
		EndPercentage: tps.EndPercentage,
	}
	waitGroup := &sync.WaitGroup{}

	for i := 0; i < tps.Connections; i++ {
		// add threshold and tails to the params
		go connection.Start(client, tps.Routes[i%len(tps.Routes)], tps.Frequency, channels[i%len(tps.Routes)], &metrics, waitGroup)
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
