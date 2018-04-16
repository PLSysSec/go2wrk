package structs

import (
	"math"
	"math/rand"
	"sync"
	"time"
	"fmt"
)

// Bootstrap is a struct that stores the latencies of all the responses.
type Bootstrap struct {
	sync.RWMutex
	List          []int
	Converged     bool
	Samples       int
	EndPercentage float64
}

// AddResponse appends a duration to the list of metrics, then returns 'Converged'. This function blocks.
func (bootstrap *Bootstrap) AddResponse(duration int) bool {
	defer bootstrap.RUnlock()
	bootstrap.RLock()

	bootstrap.List = append(bootstrap.List, duration)
	return bootstrap.Converged
}

// Start the bootstrap loop
// Bootstrap performs the bootstrapping algorithm described here: https://en.wikipedia.org/wiki/Bootstrapping_(statistics).
func (bootstrap *Bootstrap) Start() {
	for {
        if bootstrap.Converged {
            break
        }
		bootstrap.Lock()
		bootstrap.Converged = tick(bootstrap)
		bootstrap.Unlock()
        time.Sleep(time.Second / 2.0)
	}
}

func (bootstrap *Bootstrap) Check() bool{
	return bootstrap.Converged
}

func tick(bootstrap *Bootstrap) bool {
	// only start bootstrapping after the specified number of responses
	if len(bootstrap.List) < bootstrap.Samples {
		return false
	} 

	// basic bootstrapper that returns the average response time across samples
	var mean float64
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	bootstrapList := make([]float64, bootstrap.Samples)

	for i := 0; i < bootstrap.Samples; i++ {
		bootstrapList[i] = getBootstrapMean(bootstrap.List, random)
		mean += bootstrapList[i]
	}

	mean = mean / float64(bootstrap.Samples)
	variance := calculateVariance(bootstrapList, mean)
	standardDeviation := math.Sqrt(variance)

	fmt.Println(standardDeviation)
	fmt.Println(mean)
	// You are done if the deviation is less than the specified percentage
	if standardDeviation < (bootstrap.EndPercentage * mean) {
		return true
	}
	return false
}

func getBootstrapMean(metricsList []int, random *rand.Rand) float64 {
	var mean int64
	for i := 0; i < len(metricsList); i++ {
		index := random.Intn(len(metricsList))
		mean += int64(metricsList[index])
	}
	return float64(mean) / float64(len(metricsList))
}

func calculateVariance(list []float64, mean float64) float64 {
	var variance float64
	for _, value := range list {
		variance += math.Pow(value-mean, 2)
	}
	return variance / float64(len(list))
}
