package main

import (
    "sync"
    "math/rand"
    "time"
)

func Warmup(route Route, connections, total_calls int) {
    response_channel := make(chan *Response, total_calls*2)
    wait_group := &sync.WaitGroup{}

    transport := SetTLS(false)
    for i := 0; i < connections; i++ {
        wait_group.Add(1)
        go StartClient(route, response_channel, wait_group, tps.TotalCalls, transport)
    }

    wait_group.Wait()
}


func SingleNode(tps TPSReport) []byte {
    // totalCalls*2 probably so that the channel can hold resquests+responses
    var channels []chan *Response
    for i := 0; i < len(tps.Routes); i++ {
        channels = append(channels, make(chan *Response, tps.TotalCalls*2))
    }

    transport := SetTLS(*disable_keep_alives)

    bench_time := NewTimer()
    bench_time.Reset()
    // TODO check ulimit
    wait_group := &sync.WaitGroup{}

    random := rand.New(rand.NewSource(time.Now().UnixNano()))

    for i := 0; i < tps.Connections; i++ {
        // distro
        // TODO: actually write the distro switch
        index := random.Intn(len(tps.Routes)) // Generate random index
        route := tps.Routes[index]

        wait_group.Add(1)
        go StartClient(route, channels[index], wait_group, tps.TotalCalls, transport)
    }

    wait_group.Wait()

    // TODO: if we actually use this, (a json output) we will need to save an array of them
    var result []byte
    for i := 0; i < len(tps.Routes); i++ {
        result = CalcStats(channels[i], bench_time.Duration(), tps.Routes[i].Url)
    }

    return result
}
