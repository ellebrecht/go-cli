package cli

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

// Meta encapsulates meta data of a successful cli action
type Meta struct {
	Items   []*Item
	Info    string
	RawJSON string
}

// Item encapsulates generic information about something that's been created. It's normally a name / uuid pair
type Item struct {
	ID   string
	Name string
	Info string
}

// ItemFromRawJSON returns item from raw json body
func (m *Meta) ItemFromRawJSON(idKey string, nameKey string, infoKey string) (*Item, error) {
	if len(m.RawJSON) == 0 {
		return nil, errors.New("json string is empty")
	}
	dic, err := m.toMap(m.RawJSON)
	if err != nil {
		return nil, err
	}
	return m.newItem(dic, idKey, nameKey, infoKey), nil
}

// ItemsFromRawJSON returns items from raw json body
func (m *Meta) ItemsFromRawJSON(idKey string, nameKey string, infoKey string) ([]*Item, error) {
	if len(m.RawJSON) == 0 {
		return nil, errors.New("json string is empty")
	}

	array, err := m.toArray(m.RawJSON)
	if err != nil {
		return nil, err
	}
	items := []*Item{}
	for _, item := range array {
		items = append(items, m.newItem(item, idKey, nameKey, infoKey))
	}
	return items, nil
}

// UnmarshalRawJSON tries to unmarshall the raw json into the provided object
func (m *Meta) UnmarshalRawJSON(obj interface{}) error {
	return json.Unmarshal([]byte(m.RawJSON), obj)
}

// UnmarshalRawJSONAtPath tries to unmarshall the raw json path into the provided object
func (m *Meta) UnmarshalRawJSONAtPath(obj interface{}, path string) error {
	pathObj, err := m.ObjectForPath(path)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(pathObj)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, obj)
}

// ObjectForPath tries to return an object from raw json specified by the path (e.g. path.to.an.object.in.a.map)
func (m *Meta) ObjectForPath(path string) (interface{}, error) {
	var obj interface{}
	obj = m.RawJSON
	items := strings.Split(path, ".")
	for _, item := range items {
		dic, err := m.toMap(obj)
		if err != nil {
			return nil, err
		}
		obj = dic[item]
	}
	return obj, nil
}

// - private

func (m *Meta) newItem(obj map[string]interface{}, idKey string, nameKey string, infoKey string) *Item {
	item := &Item{}
	id := obj[idKey]
	name := obj[nameKey]
	info := obj[infoKey]
	if id != nil && reflect.TypeOf(id).Kind() == reflect.String {
		item.ID = reflect.ValueOf(id).String()
	}
	if name != nil && reflect.TypeOf(name).Kind() == reflect.String {
		item.Name = reflect.ValueOf(name).String()
	}
	if info != nil && reflect.TypeOf(info).Kind() == reflect.String {
		item.Info = reflect.ValueOf(info).String()
	}
	return item
}

func (m *Meta) toArray(obj interface{}) ([]map[string]interface{}, error) {
	if reflect.TypeOf(obj).Kind() != reflect.String {
		return nil, errors.New("expected string")
	}
	bytes := []byte(reflect.ValueOf(obj).String())
	array := make([]map[string]interface{}, 0)
	err := json.Unmarshal(bytes, &array)
	return array, err
}

func (m *Meta) toMap(obj interface{}) (map[string]interface{}, error) {
	if reflect.TypeOf(obj).Kind() == reflect.Map {
		return obj.(map[string]interface{}), nil
	}
	if reflect.TypeOf(obj).Kind() != reflect.String {
		return nil, errors.New("expected string")
	}
	bytes := []byte(reflect.ValueOf(obj).String())
	dic := make(map[string]interface{}, 0)
	err := json.Unmarshal(bytes, &dic)
	return dic, err
}
