package main

import (
	"github.com/kpister/go2wrk/node"
	"github.com/kpister/go2wrk/structs"
	"github.com/kpister/go2wrk/connection"

	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	tps               structs.TPSReport
	configFile        = flag.String("f", "routes.json", "the file to read routes from")
	outputDirectory   = flag.String("o", "", "the output directory to work with")
	help              = flag.Bool("h", false, "for usage")
)

func init() {
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}
	readConfig(*configFile)
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
	fmt.Println("Warming up cache on route " + tps.Routes[0].Url)
	node.Warmup(tps, 0)
	fmt.Println("Warmup complete")
    connection.Init(tps)
	node.Run(tps, *outputDirectory, 0)
}
