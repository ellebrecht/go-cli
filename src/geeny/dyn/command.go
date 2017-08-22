package dyn

import (
	"fmt"
	"net/url"

	log "geeny/log"
	"geeny/util"

	"github.com/go-openapi/spec"
)

type DynCommand struct {
	Method      string
	PathPattern string
	Params      map[string]DynParam
	Description string
	CmdPath     []string
	Hidden      bool
	Op          *spec.Operation
	Definition  string
}

type DynParam struct {
	Name        string
	In          string
	Required    bool
	Description string
	Type        string
	Flag        string
	Aliases     []string
	Hidden      bool
}

func (d DynCommand) String() string {
	return fmt.Sprintf("%v", d.PathPattern)
}

func NewCommand(op *spec.Operation, method string, path string, desc string, cmdPath []string, u *url.URL) *DynCommand {
	parameters := params(op)
	hidden := cmdMeta(op.VendorExtensible)
	if util.Contains(op.Tags, "hidden") {
		hidden = true
	}
	return &DynCommand{
		Op:          op,
		Params:      parameters,
		Method:      method,
		PathPattern: path,
		Description: desc,
		CmdPath:     cmdPath,
		Hidden:      hidden,
		Definition:  fmt.Sprintf("%s (%v %s %v) '%s' %s", responseDesc(op), op.ID, method, u, op.Description, paramsDesc(parameters)),
	}
}

func NewParam(name string, in string, required bool, desc string, typ string, flag string, aliases []string) *DynParam {
	return &DynParam{
		Name:        name,
		In:          in,
		Required:    required,
		Description: desc,
		Type:        typ,
		Flag:        flag,
		Aliases:     aliases,
	}
}

// parse operation parameters
func params(op *spec.Operation) (params map[string]DynParam) {
	params = make(map[string]DynParam)
	for _, param := range op.Parameters {
		if param.In != "body" {
			flag, aliases, isHidden := paramMeta(param.VendorExtensible, op.ID, param.Name)
			if !isHidden {
				params[param.Name] = *NewParam(param.Name, param.In, param.Required, param.Description, gettype(&param), flag, aliases)
			}
		} else if param.Schema != nil {
			for propName, propSchema := range param.Schema.Properties {
				flag, aliases, isHidden := paramMeta(propSchema.VendorExtensible, param.Schema.ID, propName)
				if !isHidden {
					params[propName] = *NewParam(propName, param.In, contains(propName, param.Schema.Required) || contains(propName, propSchema.Required), propSchema.Description, getfirsttype(&propSchema), flag, aliases)
				}
			}
		} else {
			log.Fatalf("Body param without schema: %+v", param)
		}
	}
	return
}

func cmdMeta(v spec.VendorExtensible) (hidden bool) {
	meta, ok := v.Extensions["x-meta"]
	if ok {
		log.Tracef("Meta %+v", v)
		h, ok := meta.(map[string]interface{})["hidden"]
		if ok {
			hidden = h.(bool)
		}
	}
	return
}

func paramMeta(v spec.VendorExtensible, entity string, property string) (flag string, aliases []string, isHidden bool) {
	meta, ok := v.Extensions["x-meta"]
	if !ok {
		log.Warnf("No x-meta on %v %v", entity, property)
		isHidden = true
	} else {
		f, ok := meta.(map[string]interface{})["flag"]
		if ok {
			flag = f.(string)
		} else {
			log.Warnf("No flag in x-meta on %v %v", entity, property)
		}
		al, ok := meta.(map[string]interface{})["aliases"]
		if ok {
			for _, a := range al.([]interface{}) {
				aliases = append(aliases, fmt.Sprintf("%v", a))
			}
		} else {
			log.Warnf("No aliases in x-meta on %v %v", entity, property)
		}
	}
	return
}

func gettype(param *spec.Parameter) string {
	if param != nil && param.Type != "" {
		return param.Type
	}
	return getfirsttype(param.Schema)
}

func getfirsttype(schema *spec.Schema) string {
	if schema == nil || len(schema.Type) < 0 {
		return ""
	}
	if len(schema.Type) > 1 {
		log.Warnf("Schema contains more than one type: %+v", schema)
	}
	return fmt.Sprintf("%v", schema.Type[0])
}

func contains(s string, a []string) bool {
	for _, e := range a {
		if s == e {
			return true
		}
	}
	return false
}

// gives a string representation of an operation's parameters
func paramsDesc(parameters map[string]DynParam) (params string) {
	params = ""
	i := 0
	for _, param := range parameters {
		i++
		params += "\n\t\t\t"
		params += fmt.Sprintf("%2d", i)
		params += ". "
		params += param.Name
		params += " ("
		params += param.In
		params += " "
		if param.Required {
			params += "required"
		} else {
			params += "optional"
		}
		params += " "
		params += param.Type
		params += "): '"
		params += param.Description
		params += "'"
	}
	return
}

// give a string representation of an operation's result type
func responseDesc(op *spec.Operation) (resp string) {
	resp = "[]"
	if op.Responses != nil {
		if op.Responses.Default != nil && op.Responses.Default.Schema != nil {
			resp = fmt.Sprintf("%v", op.Responses.Default.Schema.Type)
		} else if _, ok := op.Responses.StatusCodeResponses[200]; ok && op.Responses.StatusCodeResponses[200].Schema != nil {
			resp = fmt.Sprintf("%v", op.Responses.StatusCodeResponses[200].Schema.Type)
		}
	}
	return
}
