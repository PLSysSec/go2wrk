# go2wrk

A simpler, more meaninful benchmarking app modeled after go-wrk and wrk. This project is designed to stress test web apps in ways similar to organic internet traffic.

The primary addition go2wrk features is multi-route targeting. Users rarely only hit a single route of a web app while browsing, and many apps have to simultaneously handle diverse requests all competing for the same system resources. Therefore, gauging server and application performance on repeated queries to a single route in isolation is only serving a delusion. Unfortunately, this is the limited functionality of most benchmarking tools today. go2wrk benchmarks how your app and hardware perform in the presence of different usage patterns, and allows you to identify critical problems in cache and garbage collection.

go2wrk sends requests according to a given probability distribution. The simplest choice is a fair-share model where each route will be hit the same number of times on average.

### TODO

* - [ ] Build the distro model
    * - [ ] Decide what distributions we want
    * - [ ] https://en.wikipedia.org/wiki/Traffic_generation_model#Poisson_traffic_model
* - [ ] IDs in the Headers
    * - [ ] Create ids in the headers
    * - [ ] Figure out how we are using these on the end game
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
./go2wrk [flags] routes.json
```
The design philosophy we follow is that any configuration pertaining to the multiple routes should be configured from within the json file, while the app-wide features should be flags. 

A full list of flags:
```
  -CA string
    	A PEM eoncoded CA's certificate file. (default "someCertCAFile")
  -c int
    	the max numbers of connections used
  -s float
        the amount of time you want to test for (in seconds)
  -cert string
    	A PEM eoncoded certificate file. (default "someCertFile")
  -d string
    	the distribution to hit different routes
  -f string
    	json config file
  -h	for usage
  -i	TLS checks are disabled (default true)
  -k	if keep-alives are disabled (default true)
  -key string
    	A PEM encoded private key file. (default "someKeyFile")
  -t int
    	the numbers of threads used
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
        },
        {
            "Url": "https://route.com/b"
        }
    ],
    "TestTime": 10.0,
    "Threads": 8,
    "Connections": 30,
    "Distro": "Coin"
} 
```
