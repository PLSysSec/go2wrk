package connection

import (
    "github.com/kpister/go2wrk/structs"
    "github.com/kpister/go2wrk/stats"

    "io/ioutil"
    "math/rand"
    "net/http"
    "strings"
    "time"
    "sync"
    "fmt"
)

// This starts a single connection to the app. It will hit multiple routes randomly
func Start(tps structs.TPSReport, response_channels []chan *structs.Response, 
            connection_start time.Time, metrics *structs.Bootstrap, wait_group *sync.WaitGroup) {
    defer wait_group.Done()
    done := false

    random := rand.New(rand.NewSource(time.Now().UnixNano()))

    // Send a request every 1/Frequency seconds (at fastest)
    ticker := time.NewTicker(time.Second / time.Duration(tps.Frequency))
    for range ticker.C {
        // if boot channel was closed, it's time to break
        //_, ok := <-(*boot_channel)
        if done {
            ticker.Stop()
            break
        }   

        index := random.Intn(len(tps.Routes)) // Generate random index
        route := tps.Routes[index]

        request := create_request(route)

        request_start := time.Now()
        http_response, err := tps.Transport.RoundTrip(request)
        response := handle_response(http_response, err != nil)

        // hit all the described dependencies in routes.json
        for _, dependency := range route.MandatoryDependencies {
            request := create_request(dependency)
            http_response, err := tps.Transport.RoundTrip(request)
            handle_response(http_response, err != nil)
        }

        response.Duration = time.Since(request_start).Seconds()

        select {
        case response_channels[index] <- response:
            metrics.List = append(metrics.List, response.Duration)
            if len(metrics.List) > tps.Samples {
                // add response metric to bootstrap list and bootstrap
                // TODO: user specify #samples they want
                done = stats.Bootstrap(metrics, tps.Samples)
            }

            fmt.Printf("Sending requests: %.2f seconds\r", time.Since(connection_start).Seconds())
        default:
            done = true
        }
    }

}

// Used to warm up a route before we start recording results
func Warmup(tps structs.TPSReport, connection_start time.Time, wait_group *sync.WaitGroup) {
    defer wait_group.Done()

    // Send a request every 1/Frequency seconds (at fastest)
    ticker := time.NewTicker(time.Second / time.Duration(tps.Frequency))
    for range ticker.C {
        route := tps.Routes[0]
        request := create_request(route) // warmup the first route
        http_response, err := tps.Transport.RoundTrip(request)
        handle_response(http_response, err != nil)

        // hit all the described dependencies in routes.json
        for _, dependency := range route.MandatoryDependencies {
            request := create_request(dependency)
            http_response, err := tps.Transport.RoundTrip(request)
            handle_response(http_response, err != nil)
        }

        // warmups run for a set period of time (different from normal benchmarking)
        if time.Since(connection_start).Seconds() > tps.TestTime {
            ticker.Stop() 
            break
        }

        fmt.Printf("Sending requests: %.2f seconds\r", time.Since(connection_start).Seconds())
    }
}

// HELPER FUNCTIONS

// Parses a request object and prepares it to be sent
func create_request(route structs.Route) *http.Request {
    request_body_reader := strings.NewReader(route.RequestBody)
    request, _ := http.NewRequest(route.Method, route.Url, request_body_reader)

    // Split incoming header string by \n and build header pairs
    // TODO: Add counter increment to header
    header_pairs := strings.Split(route.Headers, "\n")
    for i := range header_pairs {
        split := strings.SplitN(header_pairs[i], ":", 2)
        if len(split) == 2 {
            request.Header.Set(split[0], split[1])
        }
    }
    return request
}

// Parses the response and returns it to caller
func handle_response(http_response *http.Response, err bool) *structs.Response {
    response := &structs.Response{}
    if err {
        response.Error = true
    } else {
        if http_response.ContentLength < 0 { // -1 if the length is unknown
            content, err := ioutil.ReadAll(http_response.Body)
            if err == nil {
                response.Size = int64(len(content))
            }
        } else {
            response.Size = http_response.ContentLength
        }
        response.StatusCode = http_response.StatusCode
        defer http_response.Body.Close()
    }

    return response
}