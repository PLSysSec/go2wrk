package structs

import ()

type Route struct {
	Url, Method, Headers, RequestBody string
	Dependencies []Route
}