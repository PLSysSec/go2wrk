package connection

import (
	"github.com/kpister/go2wrk/structs"

	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Start starts a single connection to the app. It will hit multiple routes randomly
// needs to take threshold and tails
func Start(route structs.Route, freq float64, responseChannel chan *structs.Response,
	connectionStart time.Time, metrics *structs.Bootstrap, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	done := false

	ticker := time.NewTicker(time.Second / time.Duration(freq))
	for range ticker.C {
		if done {
			ticker.Stop()
			break
		}

		requestStart := time.Now()
        http.Get(route.Url)
        duration := time.Since(requestStart).Nanoseconds()
		response := &structs.Response{
			Start: requestStart,
			Duration: duration,
		}

		select {
		case responseChannel <- response:
			if metrics != nil {
				done = metrics.AddResponse(response.Duration)
			}
		default:
			done = true
		}
	}
}

// Init will calibrate the app's timer
func Init(tps structs.TPSReport) {
	route := structs.Route{Url: tps.InitRoute}
	tps.Transport.RoundTrip(createRequest(route))
}

// HELPER FUNCTIONS

// Parses a request object and prepares it to be sent
func createRequest(route structs.Route) *http.Request {
	requestBodyReader := strings.NewReader(route.RequestBody)
	request, _ := http.NewRequest(route.Method, route.Url, requestBodyReader)

	// Split incoming header string by \n and build header pairs
	headerPairs := strings.Split(route.Headers, "\n")
	for i := range headerPairs {
		split := strings.SplitN(headerPairs[i], ":", 2)
		if len(split) == 2 {
			request.Header.Set(split[0], split[1])
		}
	}
	request.Header.Set("go_time", strconv.FormatInt(time.Now().UnixNano()/1000, 10))
	return request
}