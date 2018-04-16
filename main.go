package main

import (
	"github.com/kpister/go2wrk/node"
	"github.com/kpister/go2wrk/structs"
	"github.com/kpister/go2wrk/logger"

	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"strconv"
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
	tps.Logger = &logger.Logger{}
	tps.Logger.Initialize(tps.Connections * 2)
	go tps.Logger.Log()

	tps.Logger.Queue("Warming up cache on route " + tps.Routes[0].Url)
	tailThreshold, sigma := node.Warmup(tps, 0)
	tps.Logger.Queue("Warmup complete")
	tps.Logger.Queue("Threshold: " + strconv.Itoa(tailThreshold) + "\tSTD: " + strconv.Itoa(sigma))
	tps.Logger.Counter = 0
	tps.Logger.Responses = 0
    //connection.Init(tps)
	node.Barrage(tps, *outputDirectory, 0)
}

// From peterSO on Stack Exchange
func setRLimit(){
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Getting Rlimit ", err)
		os.Exit(1)
	}
	rLimit.Max = 999999
	rLimit.Cur = 999999
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println("Error Setting Rlimit ", err)
		os.Exit(1)
	}
}