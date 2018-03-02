package connection

import (
	"github.com/kpister/go2wrk/structs"

	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Start starts a single connection to the app. It will hit multiple routes randomly
func Start(tps structs.TPSReport, responseChannels []chan *structs.Response,
	connectionStart time.Time, metrics *structs.Bootstrap, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	done := false

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Send a request every 1/Frequency seconds (at fastest)
	ticker := time.NewTicker(time.Second / time.Duration(tps.Frequency))
	for range ticker.C {
		// if boot channel was closed, it's time to break
		if done {
			ticker.Stop()
			break
		}

		index := random.Intn(len(tps.Routes)) // Generate random index
		route := tps.Routes[index]

		request := createRequest(route)

		requestStart := time.Now()
		httpResponse, err := tps.Transport.RoundTrip(request)
		response := handleResponse(httpResponse, err != nil)

		// hit all the described dependencies in routes.json
		for _, dependency := range route.MandatoryDependencies {
			request := createRequest(dependency)
			httpResponse, err := tps.Transport.RoundTrip(request)
			handleResponse(httpResponse, err != nil)
		}

		response.Start = requestStart
		response.Duration = time.Since(requestStart).Nanoseconds()

		select {
		case responseChannels[index] <- response:
			done = metrics.AddResponse(response.Duration)
			fmt.Printf("Sending requests: %.1f seconds\r", time.Since(connectionStart).Seconds())
		default:
			done = true
		}
	}
}

// Warmup is used to warm up a route before we start recording results.
func Warmup(tps structs.TPSReport, connectionStart time.Time, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	// Send a request every 1/Frequency seconds (at fastest)
	ticker := time.NewTicker(time.Second / time.Duration(tps.Frequency))
	for range ticker.C {
		route := tps.Routes[0]
		request := createRequest(route) // warmup the first route
		httpResponse, err := tps.Transport.RoundTrip(request)
		handleResponse(httpResponse, err != nil)

		// hit all the described dependencies in routes.json
		for _, dependency := range route.MandatoryDependencies {
			request := createRequest(dependency)
			httpResponse, err := tps.Transport.RoundTrip(request)
			handleResponse(httpResponse, err != nil)
		}

		// warmups run for a set period of time (different from normal benchmarking)
		if time.Since(connectionStart).Seconds() > tps.MaxTestTime {
			ticker.Stop()
			break
		}

		fmt.Printf("Sending requests: %.1f seconds\r", time.Since(connectionStart).Seconds())
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

// Parses the response and returns it to caller
func handleResponse(httpResponse *http.Response, err bool) *structs.Response {
	response := &structs.Response{}
	if err {
		response.Error = true
	} else {
		if httpResponse.ContentLength < 0 { // -1 if the length is unknown
			content, err := ioutil.ReadAll(httpResponse.Body)
			if err == nil {
				response.Size = int64(len(content))
			}
		} else {
			response.Size = httpResponse.ContentLength
		}
		response.StatusCode = httpResponse.StatusCode
		defer httpResponse.Body.Close()
	}

	return response
}
