package util

import (
	"github.com/go-openapi/runtime"
)

type MockTransport struct {
	Result interface{}
	Error  error
}

func (t MockTransport) Submit(operation *runtime.ClientOperation) (interface{}, error) {
	return t.Result, t.Error
}
