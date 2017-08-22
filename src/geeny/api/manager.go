package api

import (
	"fmt"
	"reflect"
	"time"
	"encoding/json"
	"net/http"

	"geeny/net"
	"geeny/net/request"
	"geeny/net/response"
	apiResponse "geeny/api/response"
	log "geeny/log"
)

// Manager interacts with the Geeny API
type Manager struct {
	networkManager *net.Manager
	requestBuilder *request.Builder
}

// ManagerInterface contains the Manager definitions
type ManagerInterface interface {
	Perform(req request.Interface, resp response.Interface) error
}

// NewManager creates a new Manager as a ManagerInterface
func NewManager(endpoint string, timeout time.Duration) ManagerInterface {
	return &Manager{
		networkManager: net.NewManager(
			timeout,
		),
		requestBuilder: request.NewBuilder(
			endpoint,
			map[string]string{
				"Content-Type": "application/json",
			},
		),
	}
}

// Perform performs a request and receives a response
func (m *Manager) Perform(req request.Interface, resp response.Interface) error {
	// create request
	mReq, isMultipart := req.(request.MultipartInterface)
	var httpReq *http.Request
	var err error
	if isMultipart {
		httpReq, err = m.requestBuilder.NewMultipartRequest(req.RestMethod(), req.Endpoint(), req.Params(), req.Headers(), mReq.FileName(), mReq.FilePath())
	} else {
		body, err := req.Body()
		if err != nil {
			return err
		}
		httpReq, err = m.requestBuilder.NewRequest(req.RestMethod(), req.Endpoint(), req.Params(), req.Headers(), body)
	}
	if err != nil {
		return err
	}

	// perform
	statusCode, body, err := m.networkManager.PerformRequest(httpReq)
	if err != nil {
		return err
	}

	// if we got an error, try to parse the error object
	if statusCode < 200 || statusCode >= 400 {
		return m.errorResponse(body)
	}

	// marshall
	obj := resp.PointOfUnmarshall()
	if obj != nil {
		err = json.Unmarshal(body, obj)
		if err != nil {
			return err
		}
	}

	// validate
	if !resp.IsValid() {
		log.Error(reflect.TypeOf(resp).String(), "resp.IsValid() false - potential problem with api response")
	}

	return nil
}

// - private

func (m *Manager) errorResponse(body []byte) error {
	resp := new(apiResponse.ErrorResponse)
	err := json.Unmarshal(body, resp.PointOfUnmarshall())
	if err != nil {
		return err
	}
	return fmt.Errorf("error code: %d, error: %s, message: %s, validation: %v", resp.StatusCode, resp.Error, resp.Message, resp.Validation)
}
