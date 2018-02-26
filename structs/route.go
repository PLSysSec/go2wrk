package structs

type Route struct {
	Url, Method, Headers, RequestBody string
	MandatoryDependencies             []Route
	LikelyDependencies                []Route
	Samples                           int
}
