package connection

import (
    "github.com/kpister/go2wrk/structs"
    "github.com/kpister/go2wrk/stats"

    "io/ioutil"
    "math/rand"
    "net/http"
    "strings"
    "time"
    "fmt"
)

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


func Start(tps structs.TPSReport, response_channels []chan *structs.Response, connection_start time.Time, response_bootstrap *structs.Bootstrap, boot_channel *chan float64) {
    random := rand.New(rand.NewSource(time.Now().UnixNano()))

    ticker := time.NewTicker(time.Second / time.Duration(tps.Frequency))
    for range ticker.C {
        index := random.Intn(len(tps.Routes)) // Generate random index
        route := tps.Routes[index]

        request := create_request(route)

        request_start := time.Now()
        http_response, err := tps.Transport.RoundTrip(request)
        response := handle_response(http_response, err != nil)

        for _, dependency := range route.MandatoryDependencies {
            request := create_request(dependency)
            http_response, err := tps.Transport.RoundTrip(request)
            handle_response(http_response, err != nil)
        }

        response.Duration = time.Since(request_start).Seconds()

        
        //if time.Since(connection_start).Seconds() > tps.TestTime {
        //    ticker.Stop() 
        //   break
        //}
        
        // if boot channel was closed, it's time to break
        _, ok := <-(*boot_channel)
        if !ok{
            ticker.Stop()
            break
        }   

        select {
        case response_channels[index] <- response:
            // TODO: probably should only start bootstrapping after a significant # of responses are received
            // add response metric to bootstrap list and bootstrap
            response_bootstrap.Lock()
            response_bootstrap.MetricList = append(response_bootstrap.MetricList, response.Duration)
            // TODO: user specify #samples they want
            go stats.Bootstrap(response_bootstrap.MetricList, 100, boot_channel)
            response_bootstrap.Unlock()

            fmt.Printf("Sending requests: %.2f seconds\r", time.Since(connection_start).Seconds())
        }
    }

}
