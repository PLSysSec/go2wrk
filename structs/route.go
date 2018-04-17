package structs

// Route is a struct that contains information about a given route to hit on the server.
type Route struct {
	Url, Method, Headers, RequestBody string
	MandatoryDependencies             []Route
	LikelyDependencies                []Route
	Samples                           int
	Threshold                         int
}
