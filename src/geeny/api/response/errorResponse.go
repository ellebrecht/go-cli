package response

// ErrorResponse encapsulates information received from the api
type ErrorResponse struct {
	StatusCode int        `json:"statusCode"`
	Error      string     `json:"error"`
	Message    string     `json:"message"`
	Validation Validation `json:"validation"`
}

// Validation encapsulates information received from the api
type Validation struct {
	Source string   `json:"source"`
	Keys   []string `json:"keys"`
}

// - res.Interface

func (r *ErrorResponse) PointOfUnmarshall() interface{} {
	return r
}

// IsValid validates the data structure
func (r *ErrorResponse) IsValid() bool {
	return (r.StatusCode > 0 &&
		len(r.Error) > 0)
}
