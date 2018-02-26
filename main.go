package main

import (
	"github.com/kpister/go2wrk/https"
	"github.com/kpister/go2wrk/node"
	"github.com/kpister/go2wrk/structs"

	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	tps         structs.TPSReport
	connections = flag.Int("c", 0, "the max numbers of connections used")
	samples     = flag.Int("s", 0, "the max numbers of connections used")

	config_file         = flag.String("f", "routes.json", "the file to read routes from")
	output_dir          = flag.String("o", "", "the output directory to work with")
	test_time           = flag.Float64("t", 0.0, "the total runtime of the test calls")
	disable_keep_alives = flag.Bool("k", true, "if keep-alives are disabled")
	cert_file           = flag.String("cert", "someCertFile", "A PEM eoncoded certificate file.")
	key_file            = flag.String("key", "someKeyFile", "A PEM encoded private key file.")
	ca_file             = flag.String("CA", "someCertCAFile", "A PEM eoncoded CA's certificate file.")
	insecure            = flag.Bool("i", true, "TLS checks are disabled")
	help                = flag.Bool("h", false, "for usage")
)

func init() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}
	read_config(*config_file)
	initialize_tps()
}

func initialize_tps() {
	if *samples != 0 {
		tps.Samples = *samples
	}
	if *connections != 0 {
		tps.Connections = *connections
	}

	//tps.Frequency = 4
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
		Routes:      append(make([]structs.Route, 0), tps.Routes[0]),
		Connections: 10,
		TestTime:    2.0,
		Frequency:   4,
		Transport:   https.SetTLS(false, *insecure, *cert_file, *key_file, *ca_file),
	}
	node.Warmup(warmup_tps)
	fmt.Println("Warmup complete")

	fmt.Println("Starting testing")
	node.Run(tps, *output_dir)
}
