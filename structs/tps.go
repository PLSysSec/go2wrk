package structs

import "net/http"

type TPSReport struct {
    Routes []Route
    Samples int
    Connections int
    TestTime float64 
    Frequency float64
    Transport *http.Transport
}