# go2wrk

A simpler, more meaninful benchmarking app. Modeled after go-wrk and wrk, this project is designed to test web apps in ways similar to organic internet traffic.

The primary addition we feature is multi-route targetting. Users rarely only hit a single route of a web app, therefore maximizing performance of one page in isolation is only serving a delusion. By benchmarking how your app does in the presence of different usage, you can identify critical problems in cache and garbage collection.

We will send requests according to a given probability distribution. The simplest choice is a fair share model where each route will be hit on average the same number of times.

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
./go2wrk [flags] tps.json
```
The design philosophy we follow is that any configuration pertaining to the multiple routes should be configured from within the json file, while the app-wide features should be flags. 

A full list of flags:
```
  -CA string
    	A PEM eoncoded CA's certificate file. (default "someCertCAFile")
  -c int
    	the max numbers of connections used
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
  -n int
    	the total number of calls processed
  -t int
    	the numbers of threads used
```

A normal routes.json file will look something like below
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
    "TotalCalls": 100,
    "Threads": 8,
    "Connections": 30,
    "Distro": "Coin"
} 
```
