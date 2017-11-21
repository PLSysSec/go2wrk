package structs

import "net/http"

type TPSReport struct {
    Routes []Route
    Threads int
    Connections int
    Distro string // This could be a binary mapping instead
    TestTime float64 
    Frequency float64
    Transport *http.Transport
}