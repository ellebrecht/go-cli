package dyn

import (
	"time"
	"net/http"
	log "geeny/log"

	"golang.org/x/net/context"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/runtime"
)

type Request struct {
	command    *DynCommand
	params     *Params
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

func NewRequest(c *DynCommand, p *Params, t time.Duration) *Request {
	return &Request{
		command: c,
		params: p,
		timeout: t,
	}
}

func (o *Request) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {
	r.SetTimeout(o.timeout)
	for k, v := range o.params.pathParams {
		r.SetPathParam(k, v)
		log.Tracef("Adding path parameter %v: %v", k, v)
	}
	for k, v := range o.params.formParams {
		r.SetFormParam(k, v...)
		log.Tracef("Adding form parameter %v: %v", k, v)
	}
	for k, v := range o.params.queryParams {
		r.SetQueryParam(k, v...)
		log.Tracef("Adding query parameter %v: %v", k, v)
	}
	for k, v := range o.params.headerParams {
		r.SetHeaderParam(k, v...)
		log.Tracef("Adding header parameter %v: %v", k, v)
	}
	for k, v := range o.params.fileParams {
		r.SetFileParam(k, v)
		log.Tracef("Adding file parameter %v: %v", k, v)
	}
	if o.params.HasBodyParam() {
		body := o.params.GetBodyParam()
		r.SetBodyParam(body)
		log.Debugf("Adding body %v", body)
	}
	return nil
}
