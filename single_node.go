package main

import (
    "time"
    "fmt"
)

func SingleNode(tps TPSReport, output bool) {
    var channels []chan *Response
    for i := 0; i < len(tps.Routes); i++ {
        channels = append(channels, make(chan *Response, int(tps.TestTime)*tps.Connections * 10))
    }

    start := time.Now()

    for i := 0; i < tps.Connections; i++ {
        go StartClient(tps, channels, start)
    }

    time.Sleep(time.Duration(int(tps.TestTime + 1)) * time.Second)
    fmt.Println()

    if output {
        duration := time.Since(start).Seconds()

        for i, route := range tps.Routes {
            CalcStats(channels[i], duration, route.Url)
        }
    }
}
