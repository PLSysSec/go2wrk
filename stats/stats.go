package stats

import (
    "github.com/kpister/go2wrk/structs"

    "encoding/json"
    "sort"
    "fmt"
)

type Stats struct {
    Url         string
    Connections int
    Threads     int
    AvgDuration float64
    Duration    float64
    Sum         float64
    Times       []float64
    Transfered  int64
    Resp200     int64
    Resp300     int64
    Resp400     int64
    Resp500     int64
    Errors      int64
}

// TODO: need to determine return type
func Bootstrap(metrics_list []float64, samples int) bool {
    //resamples := len(metrics_list)
    return true
}

func Calculate(tps structs.TPSReport, response_channel chan *structs.Response, duration float64, url string) []byte {

    stats := &Stats{
        Url:         url,
        Connections: tps.Connections,
        Threads:     tps.Threads,
        Times:       make([]float64, len(response_channel)),
        Duration:    duration, // In seconds
        AvgDuration: duration, // In seconds
    }

    i := 0
    for res := range response_channel {
        switch {
        case res.StatusCode < 200:
            // error
        case res.StatusCode < 300:
            stats.Resp200++
        case res.StatusCode < 400:
            stats.Resp300++
        case res.StatusCode < 500:
            stats.Resp400++
        case res.StatusCode < 600:
            stats.Resp500++
        }

        stats.Sum += res.Duration // this is seconds
        stats.Times[i] = res.Duration
        i++

        stats.Transfered += res.Size

        if res.Error {
            stats.Errors++
        }

        if len(response_channel) == 0 {
            break
        }
    }

    sort.Float64s(stats.Times)

    print_stats(stats)
    b, err := json.Marshal(&stats)
    if err != nil {
        fmt.Println(err)
    }
    return b
}

func print_stats(allStats *Stats) {
    sort.Float64s(allStats.Times)
    total := float64(len(allStats.Times))
    totalInt := int64(total)
    fmt.Println("==========================BENCHMARK==========================")
    fmt.Printf("URL:\t\t\t\t%s\n\n", allStats.Url)
    fmt.Printf("Used Connections:\t\t%d\n", allStats.Connections)
    fmt.Printf("Used Threads:\t\t\t%d\n", allStats.Threads)
    fmt.Printf("Total number of calls:\t\t%d\n\n", totalInt)
    fmt.Println("===========================TIMINGS===========================")
    fmt.Printf("Total time passed:\t\t%.2fs\n", allStats.AvgDuration)
    fmt.Printf("Avg time per request:\t\t%.2fms\n", allStats.Sum/total*1e3)
    fmt.Printf("Requests per second:\t\t%.2f\n", total/(allStats.AvgDuration))
    fmt.Printf("Median time per request:\t%.2fms\n", float64(allStats.Times[(totalInt-1)/2])*1e3)
    fmt.Printf("99th percentile time:\t\t%.2fms\n", float64(allStats.Times[(totalInt/100*99)])*1e3)
    fmt.Printf("Slowest time for request:\t%.2fms\n\n", float64(allStats.Times[totalInt-1]*1e3))
    fmt.Println("=============================DATA=============================")
    fmt.Printf("Total response body sizes:\t\t%d\n", allStats.Transfered)
    fmt.Printf("Avg response body per request:\t\t%.2fms\n", float64(allStats.Transfered)/total)
    tr := float64(allStats.Transfered) / (allStats.AvgDuration )
    fmt.Printf("Transfer rate per second:\t\t%.2f Byte/s (%.2f MByte/s)\n", tr, tr/1e6)
    fmt.Println("==========================RESPONSES==========================")
    fmt.Printf("20X Responses:\t\t%d\t(%.2f%%)\n", allStats.Resp200, float64(allStats.Resp200)/total*1e2)
    fmt.Printf("30X Responses:\t\t%d\t(%.2f%%)\n", allStats.Resp300, float64(allStats.Resp300)/total*1e2)
    fmt.Printf("40X Responses:\t\t%d\t(%.2f%%)\n", allStats.Resp400, float64(allStats.Resp400)/total*1e2)
    fmt.Printf("50X Responses:\t\t%d\t(%.2f%%)\n", allStats.Resp500, float64(allStats.Resp500)/total*1e2)
    fmt.Printf("Errors:\t\t\t%d\t(%.2f%%)\n", allStats.Errors, float64(allStats.Errors)/total*1e2)
}
