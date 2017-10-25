package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "runtime"
)

type TPSReport struct {
    Routes []Route
    TotalCalls int
    Threads int
    Connections int
    Distro string // This could be a binary mapping instead
}
var (
    tps TPSReport
    threads = flag.Int("t", 0, "the numbers of threads used")
    connections = flag.Int("c", 0, "the max numbers of connections used")
    total_calls = flag.Int("n", 0, "the total number of calls processed")
    distro = flag.String("d", "", "the distribution to hit different routes")

    disable_keep_alives = flag.Bool("k", true, "if keep-alives are disabled")
    config_file = flag.String("f", "", "json config file")
    cert_file = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
    key_file = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
    ca_file = flag.String("CA", "someCertCAFile", "A PEM eoncoded CA's certificate file.")
    insecure = flag.Bool("i", true, "TLS checks are disabled")
    help = flag.Bool("h", false, "for usage")
)


func init() {
    flag.Parse()

    if *help {
        flag.PrintDefaults()
        os.Exit(1)
    }
    initialize_tps()
    // TODO handle no file error
    config_file := os.Args[len(os.Args)-1]
    if config_file != "" {
        read_config(config_file)
    }
    runtime.GOMAXPROCS(tps.Threads)
}

func initialize_tps() {
    if *threads != 0 {
        tps.Threads = *threads
    }

    if *connections != 0 {
        tps.Connections = *connections
    }

    if *total_calls != 0 {
        tps.TotalCalls = *total_calls
    }

    if *distro != "" {
        tps.Distro = *distro
    }
}

func read_config(config_file string) {
    config_data, err := ioutil.ReadFile(config_file)
    if err != nil {
        fmt.Println(err)
    }
    err = json.Unmarshal(config_data, &tps)
    if err != nil {
        fmt.Println(err)
    }
}

func main() {
    // warmup cache on first route
    // TODO: may want to make this more general in case Urls[0] is not always the first one hit
    fmt.Println("Warming up cache on route " + tps.Routes[0].Url)
    Warmup(tps.Routes[0], 10, 1000)
    fmt.Println("Warmup complete")

    SingleNode(tps)
}
