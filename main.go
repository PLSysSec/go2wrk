package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "runtime"
)

type Config struct {
    Port  string
    Nodes []string
}

var (
    numThreads        = flag.Int("t", 1, "the numbers of threads used")
    method            = flag.String("m", "GET", "the http request method")
    requestBody       = flag.String("b", "", "the http requst body")
    numConnections    = flag.Int("c", 100, "the max numbers of connections used")
    totalCalls        = flag.Int("n", 1000, "the total number of calls processed")
    disableKeepAlives = flag.Bool("k", true, "if keep-alives are disabled")
    dist              = flag.String("d", "", "dist mode")
    configFile        = flag.String("f", "", "json config file")
    config            Config
    target            string
    headers           = flag.String("H", "User-Agent: go-wrk 0.1 bechmark\nContent-Type: text/html;", "the http headers sent separated by '\\n'")
    certFile          = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
    keyFile           = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
    caFile            = flag.String("CA", "someCertCAFile", "A PEM eoncoded CA's certificate file.")
    insecure          = flag.Bool("i", false, "TLS checks are disabled")
)

func init() {
    flag.Parse()
    target = os.Args[len(os.Args)-1]
    if *configFile != "" {
        readConfig()
    }
    runtime.GOMAXPROCS(*numThreads)
}

func readConfig() {
    configData, err := ioutil.ReadFile(*configFile)
    if err != nil {
        fmt.Println(err)
    }
    err = json.Unmarshal(configData, &config)
    if err != nil {
        fmt.Println(err)
    }
}

func main() {
    SingleNode(target)
}
