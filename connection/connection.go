package connection

import (
	"github.com/kpister/go2wrk/structs"

	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Start starts a single connection to the app. It will hit multiple routes randomly
// needs to take threshold and tails
func Start(client *http.Client, route structs.Route, freq float64, responseChannel chan *structs.Response,
	metrics *structs.Bootstrap, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	done := false

	//ticker := time.NewTicker(time.Second / time.Duration(freq))
	for {
		if done {
			break
		}

		requestStart := time.Now()
		res, err := client.Get(route.Url)
        duration := time.Since(requestStart).Nanoseconds()
		if err != nil{
			panic(err)
			continue
		}
		response := &structs.Response{
			Start: requestStart,
			Duration: duration,
		}
		ioutil.ReadAll(res.Body)
		res.Body.Close()

		select {
		case responseChannel <- response:
			if metrics != nil {
				done = metrics.AddResponse(response.Duration)
			}
		default:
			done = true
		}
	}
}

// Init will calibrate the app's timer
func Init(tps structs.TPSReport) {
	//requestBodyReader := strings.NewReader("")
	request, _ := http.NewRequest("Get", tps.InitRoute, nil)
	request.Header.Set("go_time", strconv.FormatInt(time.Now().UnixNano()/1000,10))
	res, _ := http.DefaultClient.Do(request)
	res.Body.Close()
}