![](logos/go2wrk_text_away2.png "go2wrk")

[![Go Report Card](https://goreportcard.com/badge/github.com/kpister/go2wrk)](https://goreportcard.com/report/github.com/kpister/go2wrk)

A simpler, more meaningful benchmarking app modeled after go-wrk and wrk. This project is designed to stress test web apps in ways similar to organic internet traffic.

The primary addition go2wrk features is multi-route targeting. Users rarely only hit a single route of a web app while browsing, and many apps have to simultaneously handle diverse requests all competing for the same system resources. Therefore, gauging server and application performance on repeated queries to a single route in isolation is only serving a delusion. Unfortunately, this is the limited functionality of most benchmarking tools today. go2wrk benchmarks how your app and hardware perform in the presence of different usage patterns, and allows you to identify critical problems in cache and garbage collection.

go2wrk sends requests according to a given probability distribution. The simplest choice is a fair-share model where each route will be hit the same number of times on average.

### TODO

* - [ ] Update Readme
* - [ ] Look into bootstrapping
    * - [ ] make it async with a bool flag if done
    * - [ ] make sure it actually is working
    * - [ ] make it a flag
* - [ ] Add graph making code to repo in clean way
* - [ ] Compare to Autocannon
* - [ ] Build/explore node apps to test on
* - [ ] look into other performance stats other than GC
    * - [ ] Cache (in some way other than `perf`)
    * - [ ] ...
* - [ ] TLS
    * - [ ] Actually get the tls stuff working -- get those certs?
    * - [ ] Add to readme the steps needed for that

### Building and Usage

```
go get github.com/kpister/go2wrk
cd go2wrk
go build        // alternatively use go install if you have set your $GOPATH
```
You now should have an executable which you can run
```
./go2wrk [flags]        // with empty args, the program defaults to routes.json
```
The design philosophy we follow is that any configuration pertaining to the multiple routes should be configured from within the json file, while the app-wide features should be flags. 

A full list of flags:
```
  -h	for usage
  -f string
        The file name for the route descriptions. A json file (default: "routes.json")
  -o string
        The output directory for the graph data (default: current directory)
  -c int
    	the max numbers of connections used
  -s int
    	the numbers of samples to bootstrap on
  -t float
        the amount of time you want to test for (in seconds)
  -cert string
    	A PEM eoncoded certificate file. (default: "someCertFile")
  -i	TLS checks are disabled (default: true)
  -k	if keep-alives are disabled (default: true)
  -key string
    	A PEM encoded private key file. (default: "someKeyFile")
  -CA string
    	A PEM eoncoded CA's certificate file. (default: "someCertCAFile")
```

A normal routes.json file will look something like the following:
``` 
{
    "Routes": 
    [
        {
            "Url": "https://route.com/a",
            "Headers": ""                           //optional
            "Method": "Get"                         //optional
            "RequestBody": ""                       //optional
            "MandatoryDependencies":[               //optional
                // More urls
            ]
        },
        {
            "Url": "https://route.com/b"
        }
    ],
    "TestTime": 10.0,
    "Connections": 30,
    "Samples":1000,
    "Latency":1000,
    "Frequency":4
} 
```

#### Tags
go, plsyssec
