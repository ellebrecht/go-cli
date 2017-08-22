package request

// Interface defines functions needed for an http request
type Interface interface {
	Endpoint() string
	RestMethod() string
	Headers() map[string]string
	Params() map[string]string
	Body() ([]byte, error)
}
