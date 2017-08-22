package model

// Error represents api error object
type Error struct {
	StatusCode int        `json:"statusCode"`
	Error      string     `json:"error"`
	Message    string     `json:"message"`
	Validation Validation `json:"validation"`
}

// - ValidationInterface

// Validation encapsulates information received from the api
type Validation struct {
	Source string   `json:"source"`
	Keys   []string `json:"keys"`
}

// IsValid validates the data structure
func (e *Error) IsValid() bool {
	return (e.StatusCode > 0 &&
		len(e.Error) > 0)
}

func (e *Error) StringValue() string {
	return e.Message
}
