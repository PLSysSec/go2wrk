package node

import (
	"github.com/kpister/go2wrk/connection"
	"github.com/kpister/go2wrk/stats"
	"github.com/kpister/go2wrk/structs"

	"net/http"
	"sync"
)

/*
// Warmup performs a short warmup on the server. 
func Warmup(tps structs.TPSReport, index int) int{
	waitGroup := &sync.WaitGroup{}
	warmupData := make(chan *structs.Response, 100 * int64(tps.Connections))

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	for i := 0; i < tps.Connections; i++ {
		go connection.Start(tps, client, tps.Routes[index], tps.Frequency, warmupData, nil, waitGroup)
		waitGroup.Add(1)
	}
	waitGroup.Wait()
	close(warmupData)
	return stats.FindThreshold(warmupData)
}
*/

// Barrage will create connections that fire requests at the server. Then it creates the output.
func ShortBarrage(tps structs.TPSReport) {
	var channels []chan *structs.Response
	for i := 0; i < len(tps.Routes); i++ {
		channels = append(channels, make(chan *structs.Response, 200*int64(tps.Connections)))
	}
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	waitGroup := &sync.WaitGroup{}

	index := 0
	for i := 0; i < tps.Connections; i++ {
		// add threshold and tails to the params
		index = i %len(tps.Routes)
		go connection.Start(tps, client, tps.Routes[index], tps.Frequency, channels[index], nil, waitGroup)
		waitGroup.Add(1)
	}
	waitGroup.Wait()
	for i, _ := range tps.Routes{
		close(channels[i])
		tps.Routes[i].Threshold = stats.FindThreshold(channels[i])
	}
}

// Barrage will create connections that fire requests at the server. Then it creates the output.
func Barrage(tps structs.TPSReport, outputDirectory string, outputIteration int) {
	var channels []chan *structs.Response
	var metrics []structs.Bootstrap
	for i := 0; i < len(tps.Routes); i++ {
		channels = append(channels, make(chan *structs.Response, int64(1000*tps.Frequency)*int64(tps.Connections)))
		metrics = append(metrics, structs.Bootstrap{
											List: make([]int, 0), 
											Converged: false, 
											Samples: tps.Samples, 
											EndPercentage: tps.EndPercentage,
										})
	}
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	waitGroup := &sync.WaitGroup{}

	index := 0
	for i := 0; i < tps.Connections; i++ {
		// add threshold and tails to the params
		index = i %len(tps.Routes)
		go connection.Start(tps, client, tps.Routes[index], tps.Frequency, channels[index], &metrics[index], waitGroup)
		go (&metrics[index]).Start()
		waitGroup.Add(1)
	}
	waitGroup.Wait()
	tps.Logger.Kill()
	tps.Logger.Queue("Saving responses to disk")

	for i, route := range tps.Routes {
		close(channels[i])
		// update export to deal with tails
		stats.Export(channels[i], i, outputIteration, route, outputDirectory)
	}
}
