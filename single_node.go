package main

import (
    "sync"
    "fmt"
    "strconv"
    "math/rand"
    "time"
)

func Warmup(route string, numConnections, totalCalls int) {
    // Do we even care about the benchmarking here?
    responseChannel := make(chan *Response, totalCalls*2)

    wg := &sync.WaitGroup{}

    for i := 0; i < numConnections; i++ {
        //distro
        wg.Add(1)
        go StartClient(
            route,
            *headers,
            *requestBody,
            *method,
            false, //disablekeepalive
            responseChannel,
            wg,
            totalCalls,
        )
    }

    wg.Wait()
}


func SingleNode(tps TPSReport) []byte {
    // totalCalls*2 probably so that the channel can hold resquests+responses
    var channels []chan *Response
    for i := 0; i < len(tps.Routes); i++ {
        channels = append(channels, make(chan *Response, tps.TotalCalls*2))
    }

    benchTime := NewTimer()
    benchTime.Reset()
    //TODO check ulimit
    wg := &sync.WaitGroup{}

    s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)

    for i := 0; i < tps.NumConnections; i++ {
        //distro
        index := r1.Intn(len(tps.Routes)) // Generate random index
        route := tps.Routes[index]

        //fmt.Println("Starting connection " + strconv.Itoa(i) + " to " + route)
        wg.Add(1)
        go StartClient(
            route,
            *headers,
            *requestBody,
            *method,
            *disableKeepAlives, // Allow reuse of TCP connection after warmup sequence
            channels[index],
            wg,
            tps.TotalCalls,
        )
    }

    wg.Wait()

    var result []byte
    for i:= 0; i < len(tps.Routes); i++ {
        result = CalcStats(
            channels[i],
            benchTime.Duration(),
            tps.Routes[i],
        )
    }

    return result
}
