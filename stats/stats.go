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

func Export(response_channel chan *structs.Response, pos int, url string, output_dir string) {
	if output_dir != "" && string(output_dir[len(output_dir)-1]) != "/" {
		output_dir += "/"
	}

	// open a file
	data_file, _ := os.Create(output_dir + "output_" + strconv.Itoa(pos) + ".data")
	defer data_file.Close()


	output := url + "\n"

	for response := range response_channel {
		// print (time, latency)
		output += strconv.FormatInt(response.Start.UnixNano()/1000, 10) + "," + strconv.FormatInt(response.Duration, 10) + "\n" // the 10 is for the base
	}
	data_file.WriteString(output)
}

// TODO: need to determine return type
func Bootstrap(metrics *structs.Bootstrap, samples int, latency float64) bool {
	defer metrics.Unlock()
	metrics.Lock()

	// basic bootstrapper that returns the average response time across samples
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	b_mean := 0.0
	bootstrap_list := make([]float64, samples)

	for i := 0; i < samples; i++ {
		bootstrap_list[i] = get_bootstrap_mean(metrics.List, random)
		b_mean += bootstrap_list[i]
	}
	b_mean = b_mean / float64(samples)

	b_variance := calculate_variance(bootstrap_list, b_mean)
	b_standard_deviation := math.Sqrt(b_variance)

	// if s_d response time below half millisecond
	fmt.Printf("standard dev: %f\n", b_standard_deviation)
	if b_standard_deviation < latency {
		return true
	}
	return false
}

func get_bootstrap_mean(metrics_list []int64, random *rand.Rand) float64 {
	var mean int64 = 0
	for i := 0; i < len(metrics_list); i++ {
		index := random.Intn(len(metrics_list))
		mean += metrics_list[index]
	}
	return float64(mean) / float64(len(metrics_list))
}

func calculate_variance(list []float64, mean float64) float64 {
	b_variance := 0.0
	for _, value := range list {
		b_variance += math.Pow(value-mean, 2)
	}
	return b_variance / float64(len(list))
}
