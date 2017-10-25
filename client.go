package main

import (
    "io/ioutil"
    "net/http"
    "strings"
    "sync"
)

func StartClient(route Route, response_channel chan *Response, wait_group *sync.WaitGroup, total_calls int, transport *http.Transport) {
    defer wait_group.Done()

    timer := NewTimer()
    for {
        request_body_reader := strings.NewReader(route.RequestBody)
        request, _ := http.NewRequest(route.Method, route.Url, request_body_reader)
        header_pairs := strings.Split(route.Headers, "\n")

        // Split incoming header string by \n and build header pairs
        // TODO: Add counter increment to header
        for i := range header_pairs {
            split := strings.SplitN(header_pairs[i], ":", 2)
            if len(split) == 2 {
                request.Header.Set(split[0], split[1])
            }
        }

        timer.Reset()

        http_response, err := transport.RoundTrip(request)
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
            http_response.Body.Close()
        }

        response.Duration = timer.Duration()

        if len(response_channel) >= total_calls{
            break
        }
        response_channel <- response
    }

}
