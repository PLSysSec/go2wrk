package connection

import (
	"github.com/kpister/go2wrk/structs"

	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"strconv"
	"sync"
	"time"
)

// Start starts a single connection to the app. It will hit multiple routes randomly
// needs to take threshold and tails
func Start(tps structs.TPSReport, client *http.Client, route structs.Route, freq float64, responseChannel chan *structs.Response,
	metrics *structs.Bootstrap, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	done := false

	//ticker := time.NewTicker(time.Second / time.Duration(freq))
	for {
		if done {
			break
		}

		requestStart := time.Now()
		res, err := client.Get(route.Url)
		duration := int(time.Since(requestStart).Nanoseconds() / 1000)
		if err != nil {
			fmt.Println("Error connecting to server. Is it turned on?")
			os.Exit(1)
		}
		response := &structs.Response{
			Start:    requestStart,
			Duration: duration,
		}
		ioutil.ReadAll(res.Body)
		res.Body.Close()

		select {
		case responseChannel <- response:
			if metrics != nil && int(duration) > route.Threshold {
				tps.Logger.Increment()
				done = metrics.AddResponse(response.Duration)
			} else if metrics == nil {
				tps.Logger.Increment()
			} else {
				done = metrics.Check()
			}
		default:
			if metrics != nil {
				fmt.Println("Channel Full??????????????")
			}
			done = true
		}
	}
}

// Init will calibrate the app's timer
func Init(tps structs.TPSReport) {
	route := structs.Route{Url: tps.InitRoute}
	tr := &http.Transport{}
	tr.RoundTrip(createRequest(route))
}

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