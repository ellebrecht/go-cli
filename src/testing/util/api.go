package util

import (
	"reflect"

	request "geeny/net/request"
	response "geeny/net/response"
)

type MockAPIManager struct {
	Payloads []interface{}
	Error    error
}

func (a *MockAPIManager) Perform(req request.Interface, resp response.Interface) error {
	if len(a.Payloads) > 0 {
		/**

		@see http://stackoverflow.com/questions/19301742/golang-interface-to-swap-two-numbers
		basically this does
		*a, *b = *b, *a
		*resp.PointOfUnmarshall(), *a.Payload = *a.Payload, *resp.PointOfUnmarshall()

		*/
		idx := 0
		payload := a.Payloads[idx]
		a.Payloads = append(a.Payloads[:idx], a.Payloads[idx+1:]...) // pop

		valA := reflect.ValueOf(resp.PointOfUnmarshall()).Elem()
		valB := reflect.ValueOf(payload).Elem()
		tmp := valA.Interface()
		valC := reflect.ValueOf(tmp)
		valA.Set(valB)
		valB.Set(valC)
	}

	return a.Error
}
