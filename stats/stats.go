package stats

import (
	"github.com/kpister/go2wrk/structs"

	"fmt"
	"os"
	"strconv"
)

// Export is a function that outputs the time of response and latency for each request.
func Export(responseChannel chan *structs.Response, pos, iter int, url string, outputDirectory string) {
	if outputDirectory != "" && string(outputDirectory[len(outputDirectory)-1]) != "/" {
		outputDirectory += "/"
	}
	filename := outputDirectory + "output_" + strconv.Itoa(iter) + "_" + strconv.Itoa(pos) + ".data"
	output, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660);
	if err != nil {
		panic(err)
	}

	output.WriteString(url + "\n")

	for response := range responseChannel {
		line := strconv.FormatInt(response.Start.UnixNano()/1000, 10) + "," + strconv.FormatInt(response.Duration, 10) + "\n" // the 10 is for the base
		output.WriteString(line)
	}
}

// Perform will take results, and find a definition of tail latency on that group. It will return the threshold and an empty bootstrap struct tails.
func FindThreshold(responseChannel chan *structs.Response) int64 {
	var latencies []int64
	var sum int64

	fmt.Println("finding stuff")
	for response := range responseChannel {
		latencies = append(latencies, response.Duration)
		sum += response.Duration
	}

	return sum / int64(len(latencies)) * 10
}