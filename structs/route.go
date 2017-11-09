package structs

import ()

type Route struct {
	Url, Method, Headers, RequestBody string
	MandatoryDependencies []Route
	LikelyDependencies []Route
}