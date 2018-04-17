package stats

import (
	"github.com/kpister/go2wrk/structs"

	"math"
	"os"
	"strconv"
)

// Export is a function that outputs the time of response and latency for each request.
func Export(responseChannel chan *structs.Response, pos, iter int, route structs.Route, outputDirectory string) {
	if outputDirectory != "" && string(outputDirectory[len(outputDirectory)-1]) != "/" {
		outputDirectory += "/"
	}
	filename := outputDirectory + "output_" + strconv.Itoa(iter) + "_" + strconv.Itoa(pos) + ".data"
	filename2 := outputDirectory + "output_" + strconv.Itoa(iter+1) + "_" + strconv.Itoa(pos) + ".data"
	os.Remove(filename)
	os.Remove(filename2)
	output, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	output2, err := os.OpenFile(filename2, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}

	output.WriteString(route.Url + "\n")
	output2.WriteString(route.Url + "\n")

	for response := range responseChannel {
		line := strconv.FormatInt(response.Start.UnixNano()/1000, 10) + "," + strconv.Itoa(response.Duration) + "\n" // the 10 is for the base
		if response.Duration > route.Threshold {
			output.WriteString(line)
		} else {
			output2.WriteString(line)
		}
	}
}

// Perform will take results, and find a definition of tail latency on that group. It will return the threshold and an empty bootstrap struct tails.
func FindThreshold(responseChannel chan *structs.Response) int {
	var latencies []int
	var sum int

	for response := range responseChannel {
		latencies = append(latencies, int(response.Duration))
		sum += int(response.Duration)
	}

	mean := sum / len(latencies)
	sigma := int(calculateSTD(latencies, mean))

	return mean + 2*sigma
}

func calculateSTD(list []int, mean int) float64 {
	var variance float64
	for _, value := range list {
		variance += math.Pow(float64(value-mean), 2)
	}
	return math.Sqrt(variance / float64(len(list)))
}
