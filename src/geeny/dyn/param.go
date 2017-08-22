package dyn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	log "geeny/log"

	"github.com/go-openapi/runtime"
)

type Params struct {
	pathParams   map[string]string
	formParams   map[string][]string
	queryParams  map[string][]string
	headerParams map[string][]string
	fileParams   map[string]runtime.NamedReadCloser
	bodyParam    *bytes.Buffer
}

func NewParams() *Params {
	p := Params{}
	p.pathParams = make(map[string]string)
	p.formParams = make(map[string][]string)
	p.queryParams = make(map[string][]string)
	p.headerParams = make(map[string][]string)
	p.fileParams = make(map[string]runtime.NamedReadCloser)
	p.bodyParam = new(bytes.Buffer)
	return &p
}

func (params Params) AddParam(p DynParam, k string, v interface{}) {
	log.Debugf("Adding %s param: %v = %v", p.In, k, v)
	switch p.In {
	case "path":
		params.pathParams[k] = convertTo(v, p.Type).(string)
	case "header":
		params.headerParams[k] = append(params.headerParams[k], convertTo(v, p.Type).(string))
	case "query":
		params.queryParams[k] = append(params.queryParams[k], convertTo(v, p.Type).(string))
	case "form":
		params.formParams[k] = append(params.formParams[k], convertTo(v, p.Type).(string))
	case "body":
		params.appendBodyParam(p.Type, k, v)
	default:
		log.Fatalf("Unknown parameter type %v", p.In)
	}
}

func (params Params) appendBodyParam(typ string, k string, v interface{}) {
	if params.bodyParam.Len() == 0 {
		params.bodyParam.WriteRune('{')
	} else {
		params.bodyParam.WriteRune(',')
	}
	params.bodyParam.WriteRune('"')
	params.bodyParam.WriteString(k)
	params.bodyParam.WriteRune('"')
	params.bodyParam.WriteRune(':')
	params.appendBodyValue(convertTo(v, typ))
}

func (params Params) appendBodyValue(v interface{}) {
	k := reflect.ValueOf(v).Kind()
	switch k {
	case reflect.String:
		params.bodyParam.WriteRune('"')
		params.bodyParam.WriteString(v.(string))
		params.bodyParam.WriteRune('"')
	case reflect.Slice, reflect.Array:
		count := 0
		obj := v.([]string)
		params.bodyParam.WriteRune('[')
		for _, e := range obj {
			params.appendBodyValue(e)
			count = count + 1
			if count < len(obj) {
				params.bodyParam.WriteRune(',')
			}
		}
		params.bodyParam.WriteRune(']')
	case reflect.Map:
		params.bodyParam.WriteRune('{')
		count := 0
		obj := v.(map[string]interface{})
		for key, value := range obj {
			params.appendBodyValue(convertTo(key, "string").(string))
			params.bodyParam.WriteRune(':')
			params.appendBodyValue(value)
			count = count + 1
			if count < len(obj) {
				params.bodyParam.WriteRune(',')
			}
		}
		params.bodyParam.WriteRune('}')
	case reflect.Chan, reflect.Ptr, reflect.Func:
		log.Fatalf("Unsupported body param value of type %v", k)
	default:
		params.bodyParam.WriteString(fmt.Sprintf("%v", v))
	}
}

func (params Params) HasBodyParam() bool {
	return params.bodyParam.Len() > 0
}

func (params Params) GetBodyParam() string {
	params.bodyParam.WriteRune('}')
	return params.bodyParam.String()
}

//TODO
func convertTo(value interface{}, typ string) interface{} {
	log.Debugf("Converting to type %v: %v", typ, value)
	switch typ {
	case "string":
		return fmt.Sprintf("%v", value)
	case "boolean":
		boolVal, err := strconv.ParseBool(fmt.Sprintf("%v", value))
		if err != nil {
			log.Error(err)
			return false
		}
		return boolVal
	case "array":
		return strings.Split(fmt.Sprintf("%v", value), ",")
	case "object":
		obj := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(fmt.Sprintf("%v", value)), &obj)
		if err != nil {
			log.Error(err)
		}
		return obj
	}
	return value
}
