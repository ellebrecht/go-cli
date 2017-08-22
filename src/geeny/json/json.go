package json

import (
	"bytes"
	"encoding/json"
)

// PrettyJSON returns pretty json from an interface
func PrettyJSON(r interface{}) (*string, error) {
	unprettyJSON, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	var prettyJSON bytes.Buffer
	indentErr := json.Indent(&prettyJSON, unprettyJSON, "", "  ")
	if indentErr != nil {
		return nil, indentErr
	}
	str := prettyJSON.String()
	return &str, nil
}
