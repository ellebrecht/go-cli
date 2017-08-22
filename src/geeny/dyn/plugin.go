package dyn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"reflect"
	"regexp"
	"strings"
	"time"

	"geeny/cli"
	"geeny/config"
	log "geeny/log"
	"geeny/netrc"
	"geeny/output"

	"path"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
	"golang.org/x/net/context"
)

type TransportInterface interface {
	Submit(operation *runtime.ClientOperation) (interface{}, error)
}

type DynPlugin struct {
	//configuration
	SwaggerValidate    bool
	CaPemFilename      string
	Proxy              bool
	Timeout            time.Duration
	DisplayJson        bool
	RawJson            bool
	ApiUrl             string
	SwaggerUrl         string
	DefaultContentType string

	//state
	registry   strfmt.Registry
	httpClient *http.Client
	transport  TransportInterface
	host       string
	basePath   string
	schemes    []string
}

func NewPlugin(url string, swaggerJsonUrl string) *DynPlugin {
	return &DynPlugin{
		SwaggerValidate:    config.CurrentExt.SwaggerValidate,
		CaPemFilename:      config.CurrentExt.CaPemFilename,
		Proxy:              config.CurrentExt.Proxy,
		Timeout:            config.CurrentExt.Timeout,
		DisplayJson:        config.CurrentExt.DisplayJson,
		RawJson:            config.CurrentExt.RawJson,
		ApiUrl:             apiUrl(swaggerJsonUrl, url),
		SwaggerUrl:         swaggerUrl(swaggerJsonUrl, url),
		DefaultContentType: config.CurrentExt.DefaultContentType,
		registry:           strfmt.NewFormats(),
		httpClient:         nil,
		transport:          nil,
		host:               "",
		basePath:           "",
		schemes:            []string{"http", "https"},
	}
}

// initialize this plugin
func (d *DynPlugin) Init(root *cli.Command) (err error) {
	spinner := output.NewSpinner()
	spinner.Text(false, "Initializing plugin dyn")
	spinner.Start()
	defer spinner.Stop(true)

	spinner.Text(false, "Initializing swagger transport")
	d.endpointFromUrl()
	err = d.initNet()
	if err != nil {
		return
	}

	spinner.Text(false, "Loading swagger spec")
	path, err := d.readSwaggerSpec()
	if err != nil {
		return
	}

	doc, err := d.loadSpec(path)
	if err != nil {
		return
	}

	if doc.Spec() != nil && doc.Spec().Host != "" && doc.Spec().BasePath != "" && len(doc.Spec().Schemes) > 0 {
		reInit := doc.Spec().Host != d.host || doc.Spec().BasePath != d.basePath || doc.Spec().Schemes[0] != d.schemes[0]
		d.host = doc.Spec().Host
		d.basePath = doc.Spec().BasePath
		d.schemes = doc.Spec().Schemes
		if reInit {
			err = d.initNet()
			if err != nil {
				return
			}
		}
	}

	spinner.Text(false, "Processing swagger spec")
	if d.SwaggerValidate {
		err = d.validateSpec(doc)
		if err != nil {
			return
		}
	}

	err = d.resolve(doc, root)

	return
}

// close this plugin
func (d *DynPlugin) Close() error {
	return nil
}

// generic action to execute the dynamic command associated with the given context
func (d *DynPlugin) dynAction(c *cli.Context) (meta *cli.Meta, err error) {
	if c == nil {
		log.Fatal("Dynamic action context is nil")
	}

	expectedArgs := len(c.Command.Options)
	if c.Count() < expectedArgs {
		return nil, errors.New("expected " + string(expectedArgs) + " arguments")
	}

	log.Tracef("Running command %v %v", c.Command.Parent.Name, c.Command.Name)
	cmd := c.Command.Extension.(DynCommand)
	log.Tracef("Running dynamic command %+v", cmd.Definition)
	params := NewParams()
	for _, opt := range c.Args {
		p := opt.Extension.(DynParam)
		//TODO parameters: required checking, defaults/allowEmpty, type/format validation
		if opt.Value != nil {
			params.AddParam(p, opt.Name, opt.StringValue())
		} else {
			log.Tracef("Skipping param %s with nil value", opt.Name)
		}
	}
	log.TraceTime("Operation setup")

	spinner := output.NewSpinner()
	spinner.Text(false, "Connecting, please wait...")
	spinner.Start()
	defer spinner.Stop(false)

	result, err := d.transport.Submit(d.clientOp(&cmd, params, d.auth()))
	log.TraceTime("Operation roundtrip")
	if err != nil {
		log.Warnf("Error %v", err)
		return
	}

	// only generate meta data when there is no error
	meta = &cli.Meta{
		RawJSON: getResult(result, true, true, nil),
	}
	spinner.Text(false, getResult(result, d.DisplayJson, d.RawJson, err))
	return
}

func (d *DynPlugin) auth() *runtime.ClientAuthInfoWriter {
	token, err := d.token()
	if err != nil {
		return nil
	}
	if token != nil {
		auth := client.BearerToken(*token)
		return &auth
	}
	return nil
}

func (d *DynPlugin) token() (*string, error) {
	if d.host != "" && d.host != "localhost" && d.host != "127.0.0.1" {
		n, err := netrc.Instance()
		if err != nil {
			return nil, err
		}
		return n.Password(d.host)
	}
	return nil, nil
}

func (d *DynPlugin) endpointFromUrl() (err error) {
	swaggerUrl, err := url.Parse(d.ApiUrl)
	if err != nil {
		log.Warnf("Invalid url %v", d.ApiUrl)
		return err
	}
	d.host = swaggerUrl.Host
	d.basePath = swaggerUrl.Path
	d.schemes = []string{swaggerUrl.Scheme}
	return
}

// initialize the network clients
func (d *DynPlugin) initNet() (err error) {
	var caPemFile *string = nil
	if d.CaPemFilename != "" {
		caPemFile = &d.CaPemFilename
	}

	d.httpClient = NewClient(caPemFile, d.Timeout, d.Proxy)
	d.transport = client.NewWithClient(d.host, d.basePath, d.schemes, d.httpClient)
	log.Tracef("Using swagger spec %v, API endpoint %v://%v%v", d.SwaggerUrl, d.schemes, d.host, d.basePath)
	return
}

// load OpenAPI spec
func (d *DynPlugin) loadSpec(i string) (doc *loads.Document, err error) {
	doc, err = loads.Spec(i)
	if err != nil {
		return nil, err
	}
	doc, err = doc.Expanded()
	if err != nil {
		return nil, err
	}
	log.TraceTime("Loading")
	return
}

// validate OpenAPI spec
func (d *DynPlugin) validateSpec(doc *loads.Document) (err error) {
	err = validate.Spec(doc, d.registry)
	log.TraceTime("Validating")
	return
}

// load the API's swagger spec
func (d *DynPlugin) readSwaggerSpec() (string, error) {
	localSwaggerFile, ok := checkLocalFile(d.SwaggerUrl)
	if ok {
		return localSwaggerFile, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", errors.New("Unable to determine your home directory")
	}
	geenyPath := usr.HomeDir + string(os.PathSeparator) + ".geeny"
	info, err := os.Stat(geenyPath)
	if err != nil {
		err = os.Mkdir(geenyPath, 0755)
		if err != nil {
			log.Tracef("Unable to create directory %v", geenyPath)
		}
	} else if !info.IsDir() {
		return "", errors.New("Not a directory: " + geenyPath)
	}
	checkPath := geenyPath + string(os.PathSeparator) + "swagger_check.json"
	path := geenyPath + string(os.PathSeparator) + d.cacheFile()

	if info, err = os.Stat(checkPath); err == nil && !info.IsDir() {
		if info, err = os.Stat(path); err == nil && !info.IsDir() && info.Size() > 0 {
			checkContent, err := ioutil.ReadFile(checkPath)
			if err == nil {
				var check = NewSwaggerCheck()
				err = json.Unmarshal(checkContent, check)
				if err == nil {
					return d.loadSwaggerSpec(path, checkPath, check)
				} else {
					log.Tracef("Unable to parse json %v", checkPath)
				}
			} else {
				log.Warnf("Unable to read %v", checkPath)
			}
		} else {
			log.Warnf("Swagger cache does not exist: %v", path)
		}
	} else {
		log.Warnf("Swagger check does not exist: %v", checkPath)
	}

	return d.loadSwaggerSpec(path, checkPath, nil)
}

func (d *DynPlugin) cacheFile() string {
	regex := regexp.MustCompile("[^a-zA-Z0-9+=#_-]+")
	return string(regex.ReplaceAll([]byte(d.SwaggerUrl), []byte("_"))) + ".json"
}

func (d *DynPlugin) loadSwaggerSpec(path string, checkPath string, check *SwaggerCheck) (string, error) {
	token, err := d.token()
	if err != nil {
		token = nil
	}
	req, err := http.NewRequest(http.MethodGet, d.SwaggerUrl, nil)
	if err != nil {
		log.Tracef("Unable to create request for %v", d.ApiUrl)
		return "", err
	}
	if token != nil {
		req.Header.Add("Authorization", "Bearer "+*token)
	}
	if check != nil {
		checked, ok := check.Checks[d.SwaggerUrl]
		if ok {
			if checked.ETag != nil && *checked.ETag != "" {
				req.Header.Add("If-None-Match", *checked.ETag)
			}
			if checked.Timestamp != nil {
				req.Header.Add("If-Modified-Since", checked.Timestamp.UTC().Format(http.TimeFormat))
			}
		} else {
			check.Checks[d.SwaggerUrl] = SwaggerCheckElement{}
		}
	} else {
		check = NewSwaggerCheck()
		check.Checks[d.SwaggerUrl] = SwaggerCheckElement{}
	}
	now := time.Now()
	checked := check.Checks[d.SwaggerUrl]
	checked.Timestamp = &now
	resp, err := d.httpClient.Do(req)
	if err != nil {
		log.Tracef("Unable to get %v due to %v", d.ApiUrl, err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 304 {
		log.Tracef("Swagger.json is unchanged")
		return path, nil
	}
	if resp.StatusCode != 200 {
		log.Tracef("Unable to get %v, got %v", d.ApiUrl, resp.StatusCode)
		return "", err
	}
	etag := resp.Header.Get("etag")
	checked.ETag = &etag
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err == nil {
		_, err := io.Copy(file, resp.Body)
		if err == nil {
			log.Tracef("Successfully downloaded swagger.json to %v", path)
			check.Checks[d.SwaggerUrl] = checked
			checkContent, err := json.Marshal(check)
			if err == nil {
				err = ioutil.WriteFile(checkPath, checkContent, 0644)
				if err == nil {
					log.Tracef("Successfully wrote swagger check to %v", checkPath)
				} else {
					log.Tracef("Unable to write swagger check to %v", checkPath)
				}
			} else {
				log.Tracef("Unable to marshal swagger check %+v", checkContent)
			}
			return path, nil
		} else {
			return "", errors.New("Unable to write swagger to " + path)
		}
	} else {
		return "", errors.New("Unable to open for writing: " + path)
	}
}

// create a ClientOperation to execute the given Command and Params
func (d *DynPlugin) clientOp(cmd *DynCommand, params *Params, auth *runtime.ClientAuthInfoWriter) *runtime.ClientOperation {
	clientOp := runtime.ClientOperation{
		ID:                 cmd.Op.ID,
		Method:             cmd.Method,
		PathPattern:        cmd.PathPattern,
		ProducesMediaTypes: []string{d.DefaultContentType},
		ConsumesMediaTypes: []string{d.DefaultContentType},
		Schemes:            d.schemes,
		Params:             NewRequest(cmd, params, d.Timeout),
		Reader:             NewReader(cmd, d.registry),
		Context:            context.TODO(),
		Client:             d.httpClient,
	}
	if auth != nil {
		clientOp.AuthInfo = *auth
	}
	return &clientOp
}

// reconstruct spec api base path
func (d *DynPlugin) baseUrl(doc *loads.Document) (baseUrl *url.URL) {
	prefix := doc.Spec().Schemes[0] + "://" + doc.Spec().Host + doc.Spec().BasePath
	baseUrl, err := url.Parse(prefix)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	return
}

// discover and resolve all commands and add to tree
func (d *DynPlugin) resolve(doc *loads.Document, tree *cli.Command) (err error) {
	baseUrl := d.baseUrl(doc)

	for path, v := range doc.Spec().Paths.Paths {
		u := baseUrl.ResolveReference(&url.URL{Path: path})
		ops := []*spec.Operation{v.Get, v.Post, v.Put, v.Patch, v.Delete, v.Head, v.Options}
		opsMethod := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions}
		for i, op := range ops {
			err = d.addOp(op, opsMethod[i], path, u, tree)
			if err != nil {
				log.Warnf("Failed to add op", err)
				return err
			}
		}
	}

	if meta, ok := doc.Spec().Extensions["x-meta"]; ok {
		if descMeta, ok := meta.(map[string]interface{}); ok {
			if descs, ok := descMeta["descriptions"]; ok {
				if desc, ok := descs.(map[string]interface{}); ok {
					for key, val := range desc {
						cmd, ok := tree.SubCommandForPath(strings.Split(key, "_"))
						if ok {
							cmd.Summary = fmt.Sprintf("%v", val)
						}
					}
				}
			}
		}
	}

	log.TraceTime("Resolving")
	return
}

// add an operation from the OpenAPI spec to the tree
func (d *DynPlugin) addOp(op *spec.Operation, method string, path string, u *url.URL, tree *cli.Command) (err error) {
	if op != nil {
		cmd := parseOp(op, method)
		l := len(cmd)
		if l < 1 || l > 5 {
			log.Fatalf("Command with illegal category depth %v", cmd)
		}
		err = d.mergeTree(tree, *NewCommand(op, method, path, description(op), cmd, u))
	}
	return
}

// merge the dynamic command into the global command tree
func (d *DynPlugin) mergeTree(root *cli.Command, dynCmd DynCommand) error {
	l := len(dynCmd.CmdPath) - 1
	if l < 0 {
		log.Fatalf("Empty command path for Swagger operation %+v", dynCmd)
	}

	var cmd *cli.Command
	var curCmd *cli.Command
	for i, p := range dynCmd.CmdPath {
		newCmd := &cli.Command{
			Name: p,
		}
		if cmd == nil {
			cmd = newCmd
		} else {
			curCmd.Commands = append(curCmd.Commands, newCmd)
			newCmd.Parent = curCmd
		}
		if i == l {
			newCmd.Action = d.dynAction
			newCmd.Extension = dynCmd
			newCmd.Hidden = dynCmd.Hidden
			newCmd.Summary = dynCmd.Description
			newCmd.NonCategory = true
			for _, param := range dynCmd.Params {
				flag := param.Flag
				if flag == "" {
					log.Warnf("No flag on parameter %v", param.Name)
					flag = string(param.Name[0])
				}
				desc := param.Description
				if desc == "" {
					log.Warnf("No description on parameter %v", param.Name)
				}

				optionType := cli.OptionTypeString
				if strings.Compare(param.Type, "boolean") == 0 {
					optionType = cli.OptionTypeBool
				} else if strings.Compare(param.Type, "number") == 0 {
					optionType = cli.OptionTypeInt
				}

				newCmd.Options = append(newCmd.Options, &cli.Option{
					Name:        param.Name,
					Type:        optionType,
					Description: desc,
					Flag:        flag,
					Aliases:     param.Aliases,
					Extension:   param,
				})
			}
		}
		curCmd = newCmd
	}
	return root.Merge(cmd)
}

func checkLocalFile(s string) (string, bool) {
	u, err := url.Parse(s)
	if err == nil && (u.Scheme == "file" || u.Scheme == "" || u.Host == "") {
		return u.Path, true
	}
	return "", false
}

type SwaggerCheck struct {
	Checks map[string]SwaggerCheckElement `json:"checks"`
}

type SwaggerCheckElement struct {
	Timestamp *time.Time `json:"timestamp"`
	ETag      *string    `json:"etag"`
}

func NewSwaggerCheck() *SwaggerCheck {
	return &SwaggerCheck{make(map[string]SwaggerCheckElement)}
}

// get a readable operation description from the swagger elements
func description(op *spec.Operation) (desc string) {
	desc = op.Description
	if desc == "" {
		desc = op.Summary
	}
	if desc == "" {
		log.Warnf("No description or summary on operation %v", op.ID)
	}
	return
}

// parse an operation and yield tag and command names
func parseOp(op *spec.Operation, method string) (cmd []string) {
	// first try to get both from parsing the operationId
	cmd, ok := parseOpId(op.ID)

	// as a fallback use the HTTP method to determine command and get tag from declared tags
	if !ok {
		if len(op.Tags) > 0 {
			cmd = append(cmd, op.Tags[0])
		}

		switch method {
		case http.MethodGet:
			if pathParams(op) == 0 {
				cmd = append(cmd, "list")
			} else {
				cmd = append(cmd, "get")
			}
		case http.MethodPost:
			cmd = append(cmd, "create")
		case http.MethodPut, http.MethodPatch:
			cmd = append(cmd, "update")
		case http.MethodDelete:
			cmd = append(cmd, "delete")
		default:
			cmd = append(cmd, method)
		}
	}
	return
}

// parse an operationId name
func parseOpId(opId string) ([]string, bool) {
	parts := strings.Split(opId, "_")
	if len(parts) <= 0 {
		log.Warnf("Cannot parse opId '%v'", opId)
		return []string{}, false
	}
	return parts, true
}

// count the number of path parameters of an operation
func pathParams(op *spec.Operation) (params int) {
	if op != nil {
		for _, param := range op.Parameters {
			if param.In == "path" {
				params++
			}
		}
	}
	return
}

// print an operation's result value
func getResult(r interface{}, json bool, rawJson bool, err error) string {
	if err != nil {
		return fmt.Sprintf("Error %v\n", err)
	}
	if r == nil {
		return "<nil>\n"
	} else {
		switch r := r.(type) {
		case Response:
			if json {
				return r.Json(rawJson)
			}
			return r.String()
		case *Response:
			if json {
				return (*r).Json(rawJson)
			}
			return (*r).String()
		default:
			return fmt.Sprintf("(%v) \"%v\"\n", reflect.TypeOf(r), r)
		}
	}
}

func swaggerUrl(swaggerUrl string, apiUrl string) string {
	if swaggerUrl != "" {
		return swaggerUrl
	}
	if strings.LastIndex(apiUrl, "/") < len(apiUrl)-1 {
		return apiUrl + "/swagger.json"
	}
	return apiUrl + "swagger.json"
}

func apiUrl(swaggerUrl string, apiUrl string) string {
	if apiUrl != "" {
		return apiUrl
	}
	u, err := url.Parse(swaggerUrl)
	if err != nil {
		log.Warnf("Error parsing url: %v", err)
		return ""
	}
	return path.Dir(u.Path)
}
