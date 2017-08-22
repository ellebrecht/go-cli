package response

// Interface defines functions needed for an http response
type Interface interface {
	PointOfUnmarshall() interface{}
	IsValid() bool
}
