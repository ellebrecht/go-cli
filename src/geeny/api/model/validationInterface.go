package model

// ValidationInterface defines functions to validate a model
type ValidationInterface interface {
	IsValid() bool
}
