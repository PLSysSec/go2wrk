package node

import (
    "github.com/kpister/go2wrk/connection"
    "github.com/kpister/go2wrk/structs"
    "github.com/kpister/go2wrk/stats"

    "strconv"
    "time"
    "sync"
    "fmt"
)

func Warmup(tps structs.TPSReport) {
    wait_group := &sync.WaitGroup{}
    start := time.Now()
    for i := 0; i < tps.Connections; i++ {
        go connection.Warmup(tps, start, wait_group)
        wait_group.Add(1)
    }
    wait_group.Wait()
    fmt.Println()
}

func Run(tps structs.TPSReport) {
    var channels []chan *structs.Response
    for i := 0; i < len(tps.Routes); i++ {
        channels = append(channels, make(chan *structs.Response, int(tps.TestTime)*tps.Connections * 10))
    }

    // shared response metric collector and corresponding lock
    metrics := structs.Bootstrap{
        List: make([]int64, 0),
    }
    wait_group := &sync.WaitGroup{}
    start := time.Now()

    for i := 0; i < tps.Connections; i++ {
        go connection.Start(tps, channels, start, &metrics, wait_group)
        wait_group.Add(1)
    }

    wait_group.Wait()
    fmt.Println()
    fmt.Print("Done with that")

    //duration := time.Since(start).Seconds()

    for i, _ := range tps.Routes {
        close(channels[i])
        //stats.Calculate(tps, channels[i], duration, route.Url)
        stats.Export(channels[i], strconv.Itoa(i))
    }
    fmt.Printf("Response numbers: %d\n", len(metrics.List))
}
