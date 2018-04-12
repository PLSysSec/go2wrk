package main

import (
	"github.com/kpister/go2wrk/node"
	"github.com/kpister/go2wrk/structs"

	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
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

	setRLimit()
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
	tailThreshold := node.Warmup(tps, 0)
	fmt.Printf("Threshold Found %d\n", tailThreshold)
	fmt.Println("Warmup complete")
    //connection.Init(tps)
	node.Barrage(tps, *outputDirectory, 0)
}

// From peterSO on Stack Exchange
func setRLimit(){
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
	}
	rLimit.Max = 999999
	rLimit.Cur = 999999
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
	}
}