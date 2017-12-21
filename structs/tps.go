package structs

import "net/http"

type TPSReport struct {
    Routes []Route
    Threads int
    Connections int
    TestTime float64 
    Frequency float64
    Transport *http.Transport
}