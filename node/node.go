package node

import (
    "github.com/kpister/go2wrk/connection"
    "github.com/kpister/go2wrk/structs"
    "github.com/kpister/go2wrk/stats"

    "time"
    "fmt"
)

func Run(tps structs.TPSReport, output bool) {
    var channels []chan *structs.Response
    for i := 0; i < len(tps.Routes); i++ {
        channels = append(channels, make(chan *structs.Response, int(tps.TestTime)*tps.Connections * 10))
    }

    // shared response metric collector and corresponding lock
    response_bootstrap := structs.Bootstrap{
        MetricList: make([]float64, 0),
    }
    boot_channel := make(chan float64)

    start := time.Now()

    for i := 0; i < tps.Connections; i++ {
        go connection.Start(tps, channels, start, &response_bootstrap, &boot_channel)
    }

    //time.Sleep(time.Duration(int(tps.TestTime + 1)) * time.Second)
    fmt.Println()

    if output {
        duration := time.Since(start).Seconds()

        for i, route := range tps.Routes {
            stats.Calculate(tps, channels[i], duration, route.Url)
        }
    }
}
