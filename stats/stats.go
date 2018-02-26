package stats

import (
	"github.com/kpister/go2wrk/structs"

	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
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

// Bootstrap performs the bootstrapping algorithm described here: https://en.wikipedia.org/wiki/Bootstrapping_(statistics).
func Bootstrap(metrics *structs.Bootstrap, samples int, latency float64) bool {
	defer metrics.Unlock()
	metrics.Lock()

	// basic bootstrapper that returns the average response time across samples
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	bootstrapMean := 0.0
	bootstrapList := make([]float64, samples)

	for i := 0; i < samples; i++ {
		bootstrapList[i] = getBootstrapMean(metrics.List, random)
		bootstrapMean += bootstrapList[i]
	}
	bootstrapMean = bootstrapMean / float64(samples)

	bootstrapVariance := calculateVariance(bootstrapList, bootstrapMean)
	bootstrapStandardDeviation := math.Sqrt(bootstrapVariance)

	// if s_d response time below half millisecond
	fmt.Printf("standard dev: %f\n", bootstrapStandardDeviation)
	if bootstrapStandardDeviation < latency {
                return true // TODO: set a bootstrap struct var to true -- this will determine if done for everyone
	}
	return false
}

func getBootstrapMean(metricsList []int64, random *rand.Rand) float64 {
	var mean int64 = 0
	for i := 0; i < len(metricsList); i++ {
		index := random.Intn(len(metricsList))
		mean += metricsList[index]
	}
	return float64(mean) / float64(len(metricsList))
}

func calculateVariance(list []float64, mean float64) float64 {
	bootstrapVariance := 0.0
	for _, value := range list {
		bootstrapVariance += math.Pow(value-mean, 2)
	}
	return bootstrapVariance / float64(len(list))
}
