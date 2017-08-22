package dyn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	log "geeny/log"
	"geeny/output"
	"geeny/util"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
)

type Response struct {
	runtime.ClientResponse
	Payload interface{}
	Body    string
	Display string
	command *DynCommand
}

func NewResponse(c *DynCommand) *Response {
	return &Response{
		command: c,
	}
}

func (o *Response) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) (err error) {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(response.Body(), buf)
	err = consumer.Consume(tee, &o.Payload)
	o.Body = buf.String()
	log.Tracef("\tBody: %v", o.Body)
	if err != nil && err != io.EOF {
		return
	}
	o.ClientResponse = response
	o.typedResponse()
	return nil
}

func (o *Response) typedResponse() {
	if o.ClientResponse == nil {
		log.Tracef("Empty response %+v", o)
		return
	}

	code := o.Code()
	var def bool
	var rt spec.Response
	var ok bool
	payloadType := reflect.TypeOf(o.Payload)
	if o.command.Op.Responses.StatusCodeResponses != nil {
		rt, ok = o.command.Op.Responses.StatusCodeResponses[code]
		def = false
	}
	if !ok && o.command.Op.Responses.Default != nil {
		rt = *o.command.Op.Responses.Default
		ok = true
		def = true
	}

	if !ok {
		log.Warnf("No response definition for %v and no default, got type %v", code, payloadType)
		return
	}
	if rt.Schema == nil {
		log.Warnf("No schema (rt.Schema) %v", rt)
		return
	}

	log.Tracef("Expected %v (default %v) response of type %v, got %v", code, def, rt.Schema.Type, payloadType)
	rtst := rt.Schema.Type
	if len(rtst) <= 0 {
		log.Warnf("Response definition for %v contains no elements (%v), got type %v", code, rtst, payloadType)
		o.Display = stringDisplay(o.Payload)
	} else {
		if len(rtst) > 1 {
			log.Warnf("Response definition for %v contains more than one element (%v), got type %v", code, rtst, payloadType)
		}
		dispType, dispArgs, ok := displayMeta(rt.VendorExtensible)
		if !ok || dispType == "" {
			log.Warnf("No resultDisplay definition for %v", code, rtst, payloadType)
			dispType = "json"
		}

		var tableData []interface{}
		switch rtst[0] {
		case "null":
			if o.Payload != nil {
				log.Warnf("Expected %v (default %v) response of type %v, got %v", code, def, rtst[0], payloadType)
			}
		case "boolean":
			_, ok := o.Payload.(bool)
			if !ok {
				log.Warnf("Expected %v (default %v) response of type %v, got %v", code, def, rtst[0], payloadType)
			}
		case "integer":
			_, ok := o.Payload.(int)
			if !ok {
				log.Warnf("Expected %v (default %v) response of type %v, got %v", code, def, rtst[0], payloadType)
			}
		case "number":
			_, ok := o.Payload.(float32)
			if !ok {
				_, ok := o.Payload.(float64)
				if !ok {
					log.Warnf("Expected %v (default %v) response of type %v, got %v", code, def, rtst[0], payloadType)
				}
			}
		case "string":
			_, ok := o.Payload.(string)
			if !ok {
				log.Warnf("Expected %v (default %v) response of type %v, got %v", code, def, rtst[0], payloadType)
			}
		case "object":
			listData, ok := o.Payload.(map[string]interface{})
			if !ok {
				log.Warnf("Expected %v (default %v) response of type %v, got %v", code, def, rtst[0], payloadType)
			} else {
				tableData = []interface{}{listData}
			}
		case "array", "slice":
			tableData, ok = o.Payload.([]interface{})
			if !ok {
				log.Warnf("Expected %v (default %v) response of type %v, got %v", code, def, rtst[0], payloadType)
			}
		default:
			log.Warnf("No valid response type definition for %v and no default, got type %v", code, payloadType)
		}

		switch dispType {
		case "message":
			o.Display = messageDisplay(o.Payload, dispArgs)
		case "list":
			o.Display = listDisplay(tableData, dispArgs)
		case "table":
			o.Display = tableDisplay(tableData, dispArgs)
		case "json":
			o.Display = o.Json(false)
		default:
			log.Warnf("Unsupported display type %v", dispType)
			o.Display = stringDisplay(o.Payload)
		}
	}
}

func displayMeta(rt spec.VendorExtensible) (typ string, args map[string]interface{}, ok bool) {
	meta, ok := rt.Extensions["x-meta"]
	if !ok || meta == nil {
		return
	}
	metaMap, ok := meta.(map[string]interface{})
	if !ok || metaMap == nil {
		return
	}
	rd, ok := metaMap["resultDisplay"]
	if !ok || rd == nil {
		return
	}
	rdMap, ok := rd.(map[string]interface{})
	if !ok || rdMap == nil {
		return
	}
	rdType, ok := rdMap["type"]
	if !ok || rdType == nil {
		return
	}
	rdTypeString, ok := rdType.(string)
	if !ok {
		return
	}
	return rdTypeString, rdMap, true
}

func (o *Response) String() string {
	if o.Display != "" {
		return o.Display
	}
	return o.Json(false)
}

func (o *Response) Json(raw bool) string {
	if raw {
		return o.Body
	}
	return jsonDisplay(o.Payload, 0, false, "\n")
}

func stringDisplay(data interface{}) string {
	return fmt.Sprintf("%v", data)
}

func messageDisplay(data interface{}, args map[string]interface{}) (res string) {
	formati, ok := args["format"]
	if !ok {
		log.Warnf("Display type message is missing format")
		return stringDisplay(data)
	}
	format, ok := formati.(string)
	if !ok {
		log.Warnf("Display type message has invalid format")
		return stringDisplay(data)
	}
	formata, ok := args["formatArgs"]
	if ok {
		formatArgs, ok := formata.([]interface{})
		if !ok {
			log.Warnf("Display type message has invalid formatArgs")
			return stringDisplay(data)
		}
		structData, ok := data.(map[string]interface{})
		if !ok {
			log.Warnf("Display type message with incompatible data: not an object")
			return stringDisplay(data)
		}
		args := []interface{}{}
		for _, a := range formatArgs {
			args = append(args, structData[a.(string)])
		}
		return fmt.Sprintf(format, args...)
	} else if strings.Contains(format, "%") && data != nil {
		return fmt.Sprintf(format, data)
	} else {
		return format
	}
}

func fields(args map[string]interface{}) (fields []string, titles []string, ok bool) {
	fieldsi, ok := args["fields"]
	if !ok {
		return nil, nil, false
	} else {
		l, ok := fieldsi.([]interface{})
		if ok {
			fields = make([]string, len(l))
			titles = make([]string, len(l))
			for i, a := range l {
				field, ok := a.(map[string]interface{})
				if ok {
					n, ok := field["id"]
					if ok {
						t, ok := field["name"]
						if !ok {
							t = n
						}
						fields[i] = fmt.Sprintf("%v", n)
						titles[i] = fmt.Sprintf("%v", t)
					}
				}
			}
			return fields, titles, true
		}
		return nil, nil, false
		if !ok {
			return nil, nil, false
		}
	}
	return
}

func tableDisplay(data []interface{}, args map[string]interface{}) (res string) {
	if data == nil || len(data) <= 0 {
		return "(empty)"
	}
	data0, ok := data[0].(map[string]interface{})
	if !ok {
		log.Warnf("Display type table on invalid data")
		return stringDisplay(data)
	}
	fields, titles, ok := fields(args)
	if !ok {
		fields = mapKeys(data0)
		titles = fields
		log.Warnf("Display type table has invalid fields/titles")
	}
	s, _ := output.NewMapTable(titles, fields, data).String()
	return *s
}

func listDisplay(rows []interface{}, args map[string]interface{}) string {
	if rows == nil || len(rows) <= 0 {
		return "(empty)"
	}
	fields, titles, ok := fields(args)
	if !ok {
		fields = mapKeys(rows[0].(map[string]interface{}))
		titles = fields
		log.Warnf("Display type list has invalid fields/titles")
	}
	w := 0
	for _, t := range titles {
		w = util.Max(w, len(t))
	}
	var buffer bytes.Buffer
	maxRow := len(rows) - 1
	for r, row := range rows {
		data, ok := row.(map[string]interface{})
		if ok {
			for i, k := range fields {
				buffer.WriteString("\x1b[31m")
				buffer.WriteString(util.Pad(titles[i], w))
				buffer.WriteString("\x1b[0m: ")
				buffer.WriteString(fmt.Sprintf("%v", data[k]))
				buffer.WriteString("\n")
			}
			if r < maxRow {
				buffer.WriteString("\n")
			}
		} else {
			log.Warnf("Not an object: %v", row)
		}
	}
	return buffer.String()
}

func mapKeys(data map[string]interface{}) (fields []string) {
	fields = make([]string, len(data))
	i := 0
	for k, _ := range data {
		fields[i] = k
		i++
	}
	return
}

func jsonDisplay(data interface{}, i int, f bool, s string) (res string) {
	res = ""
	tab := tabs(i)
	if f {
		res += tab
	}
	if data == nil {
		res += "<nil>"
	} else {
		switch r := data.(type) {
		case []interface{}:
			res += "[\n"
			for _, e := range r {
				res += jsonDisplay(e, i+1, true, ",\n")
			}
			res += fmt.Sprintf("%s]", tab)
		case map[string]interface{}:
			res += "{\n"
			for k, v := range r {
				res += fmt.Sprintf("%s\"%v\": ", tabs(i+1), k)
				res += jsonDisplay(v, i+1, false, ",\n")
			}
			res += fmt.Sprintf("%s}", tab)
		case string:
			res += fmt.Sprintf("\"%s\"", strings.Replace(r, "\n", "\n"+tabs(i+1), -1))
		case json.Number:
			res += fmt.Sprintf("%v", r)
		default:
			res += fmt.Sprintf("(%v) %v", reflect.TypeOf(r), r)
		}
	}
	res += s
	return
}

func tabs(i int) string {
	return strings.Repeat("  ", i)
}
