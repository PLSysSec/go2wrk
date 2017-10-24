package main

import (
    "sync"
    "fmt"
    "strconv"
)

func SingleNode(toCall string, numConnections, totalCalls int, isWarmup bool) []byte {
    // totalCalls*2 probably so that the channel can hold resquests+responses
    responseChannel := make(chan *Response, totalCalls*2)

    benchTime := NewTimer()
    benchTime.Reset()
    //TODO check ulimit
    wg := &sync.WaitGroup{}

    for i := 0; i < numConnections; i++ {
        fmt.Println("Starting connection " + strconv.Itoa(i) + " to " + toCall)
        wg.Add(1)
        go StartClient(
            toCall,
            *headers,
            *requestBody,
            *method,
            *disableKeepAlives,
            responseChannel,
            wg,
            totalCalls,
        )
    }

    wg.Wait()

    result := make([]byte, 0)

    fmt.Println("MADE IT BEFORE STATISTICS")
    
    if !isWarmup {
        fmt.Println("PRINTING STATS")
        result = CalcStats(
            responseChannel,
            benchTime.Duration(),
            toCall,
        )
    }

    //result = CalcStats(
      //  responseChannel,
        //benchTime.Duration(),
        //toCall,
    //)

    return result
}
