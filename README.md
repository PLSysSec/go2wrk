# go2wrk

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
./go2wrk [flags] 
```
The design philosophy we follow is that any configuration pertaining to the multiple routes should be configured from within the json file, while the app-wide features should be flags. 

A full list of flags:
```
  -f string
        The file name for the route descriptions. A json file
  -CA string
    	A PEM eoncoded CA's certificate file. (default "someCertCAFile")
  -c int
    	the max numbers of connections used
  -t float
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


## License

This Software is licensed under the MIT License.

Copyright (c) 2013 adeven GmbH,
http://www.adeven.com

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

#### Tags
go, plsyssec
