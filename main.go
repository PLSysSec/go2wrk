package main

import (
    "github.com/kpister/go2wrk/structs"
    "github.com/kpister/go2wrk/https"
    "github.com/kpister/go2wrk/node"

    "encoding/json"
    "io/ioutil"
    "runtime"
    "flag"
    "fmt"
    "os"
    
)

var (
    tps structs.TPSReport
    threads = flag.Int("t", 0, "the numbers of threads used")
    connections = flag.Int("c", 0, "the max numbers of connections used")
    distro = flag.String("d", "", "the distribution to hit different routes")

    test_time = flag.Float64("s", 0.0, "the total runtime of the test calls")
    disable_keep_alives = flag.Bool("k", true, "if keep-alives are disabled")
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
    // TODO handle no file error
    config_file := os.Args[len(os.Args)-1]
    if config_file != "" {
        read_config(config_file)
    }

    initialize_tps()
    fmt.Println(runtime.GOMAXPROCS(tps.Threads))
    fmt.Println(runtime.NumCPU())
}

func initialize_tps() {
    if *threads != 0 { tps.Threads = *threads }
    if *connections != 0 { tps.Connections = *connections }
    if *distro != "" { tps.Distro = *distro }

    tps.Frequency = 2 
    tps.Transport = https.SetTLS(*disable_keep_alives, *insecure, *cert_file, *key_file, *ca_file)
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

    warmup_tps := structs.TPSReport{
        Routes: append(make([]structs.Route, 0), tps.Routes[0]),
        Connections: 10,
        Distro: "coin",
        TestTime: 5.0,
        Frequency: 2,
        Transport: https.SetTLS(false, *insecure, *cert_file, *key_file, *ca_file),
    }
    node.Warmup(warmup_tps)
    fmt.Println("Warmup complete")

    fmt.Println("Starting testing")
    node.Run(tps)
}
