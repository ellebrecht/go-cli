package dyn

import (
	"encoding/json"
	"fmt"
	"net/http"

	"geeny/api/model"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

type Reader struct {
	command *DynCommand
	formats strfmt.Registry
}

func NewReader(cmd *DynCommand, registry strfmt.Registry) *Reader {
	return &Reader{
		command: cmd,
		formats: registry,
	}
}

func (r *Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	// read response
	code := response.Code()
	result := NewResponse(r.command)
	if err := result.readResponse(response, consumer, r.formats); err != nil {
		return nil, err
	}

	// handle non error
	switch code {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNonAuthoritativeInfo, http.StatusNoContent, http.StatusResetContent:
		return result, nil
	}

	// handle error
	errModel := &model.Error{}
	err := json.Unmarshal([]byte(result.Body), errModel)
	if err != nil {
		return nil, err
	}
	errPreText := "There was a problem processing your request"
	if len(errModel.StringValue()) == 0 {
		return nil, fmt.Errorf("\n%s. Please try one of the following:\n%s\n%s\n%s",
			errPreText,
			"* Try again after a few moments",
			"* Update your cli to the newest version",
			"* Contact support@geeny.io")
	}
	return nil, fmt.Errorf("\n%s:\n%s", errPreText, errModel.StringValue())
}
