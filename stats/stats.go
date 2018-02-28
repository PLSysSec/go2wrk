package stats

import (
	"github.com/kpister/go2wrk/structs"

	"os"
	"strconv"
)

// Export is a function that outputs the time of response and latency for each request.
func Export(responseChannel chan *structs.Response, pos int, url string, outputDirectory string) {
	if outputDirectory != "" && string(outputDirectory[len(outputDirectory)-1]) != "/" {
		outputDirectory += "/"
	}

	// open a file
	dataFile, _ := os.Create(outputDirectory + "output_" + strconv.Itoa(pos) + ".data")
	defer dataFile.Close()

	output := url + "\n"

	for response := range responseChannel {
		// print (time, latency)
		output += strconv.FormatInt(response.Start.UnixNano()/1000, 10) + "," + strconv.FormatInt(response.Duration, 10) + "\n" // the 10 is for the base
	}
	dataFile.WriteString(output)
}