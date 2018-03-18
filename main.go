package main

import (
	"github.com/kpister/go2wrk/https"
	"github.com/kpister/go2wrk/node"
	"github.com/kpister/go2wrk/structs"
	fishStructs "github.com/streddy/go-fish/structs"

	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	tps               structs.TPSReport
	connections       = flag.Int("c", 0, "the max numbers of connections used")
	samples           = flag.Int("s", 0, "the max numbers of connections used")
	testTime          = flag.Float64("t", 0.0, "the total runtime of the test calls")
	configFile        = flag.String("f", "routes.json", "the file to read routes from")
	outputDirectory   = flag.String("o", "", "the output directory to work with")
	certFile          = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
	keyFile           = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
	caFile            = flag.String("CA", "someCertCAFile", "A PEM eoncoded CA's certificate file.")
	disableKeepAlives = flag.Bool("k", true, "if keep-alives are disabled")
	insecure          = flag.Bool("i", true, "TLS checks are disabled")
	help              = flag.Bool("h", false, "for usage")
)

func init() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}
	readConfig(*configFile)
	initializeTPS()
}

func initializeTPS() {
	if *samples != 0 {
		tps.Samples = *samples
	}
	if *connections != 0 {
		tps.Connections = *connections
	}

	//tps.Frequency = 4
	tps.Transport = https.SetTLS(*disableKeepAlives, *insecure, *certFile, *keyFile, *caFile)
}

func readConfig(configFile string) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(configData, &tps)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// warmup cache on first route
	// TODO: may want to make this more general in case Urls[0] is not always the first one hit
	fmt.Println("Warming up cache on route " + tps.Routes[0].Url)

	warmupTPS := structs.TPSReport{
		Routes:      append(make([]fishStructs.Route, 0), tps.Routes[0]),
		Connections: 10,
		MaxTestTime: 2.0,
		Frequency:   4,
		Transport:   https.SetTLS(false, *insecure, *certFile, *keyFile, *caFile),
	}
	node.Warmup(warmupTPS)
	fmt.Println("Warmup complete")

	// TODO: continue integrating here
	fmt.Println("Starting testing")
	node.Run(tps, *outputDirectory)
}
