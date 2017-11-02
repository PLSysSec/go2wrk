package main

import (
    "io/ioutil"
    "math/rand"
    "net/http"
    "strings"
    "time"
    "fmt"
)

func StartClient(tps TPSReport, response_channels []chan *Response, connection_start time.Time) {
    random := rand.New(rand.NewSource(time.Now().UnixNano()))

    ticker := time.NewTicker(time.Second / time.Duration(tps.Frequency))
    for range ticker.C {
        index := random.Intn(len(tps.Routes)) // Generate random index
        route := tps.Routes[index]

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

        request_start := time.Now()
        http_response, err := tps.Transport.RoundTrip(request)
        response := &Response{}

        if err != nil {
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
            // This will supposedly stop too many port errors
            defer http_response.Body.Close()
        }

        response.Duration = time.Since(request_start).Seconds()

        if time.Since(connection_start).Seconds() > tps.TestTime {
            ticker.Stop() 
            break
        }

        select {
        case response_channels[index] <- response:
            fmt.Printf("Sending requests: %.2f seconds\r", time.Since(connection_start).Seconds())
        }
    }

}
