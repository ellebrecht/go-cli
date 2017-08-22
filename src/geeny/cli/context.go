package cli

import (
	"errors"
	"reflect"
	"strings"
)

// Context is passed back to an action containing the command line args
type Context struct {
	Args    []*Option
	Command *Command
}

// GetString returns a string arg
func (c *Context) GetString(index int) (*string, error) {
	str := c.get(index)
	err := c.validate(c.Args[index].Name, str, reflect.String)
	if err != nil {
		return nil, err
	}
	if len(*str.(*string)) == 0 {
		return nil, errors.New(c.Args[index].Name + " is missing")
	}
	return str.(*string), nil
}

// GetStringForFlag returns a string arg
func (c *Context) GetStringForFlag(key string) (*string, error) {
	return c.GetString(c.getIndex(key))
}

// GetBool returns a string arg
func (c *Context) GetBool(index int) (*bool, error) {
	boolVal := c.get(index)
	err := c.validate(c.Args[index].Name, boolVal, reflect.Bool)
	if err != nil {
		return nil, err
	}
	return boolVal.(*bool), nil
}

// GetBoolForFlag returns a string arg
func (c *Context) GetBoolForFlag(key string) (*bool, error) {
	return c.GetBool(c.getIndex(key))
}

// GetInt returns a string arg
func (c *Context) GetInt(index int) (*int, error) {
	intVal := c.get(index)
	err := c.validate(c.Args[index].Name, intVal, reflect.Int)
	if err != nil {
		return nil, err
	}
	return intVal.(*int), nil
}

// GetIntForFlag returns a string arg
func (c *Context) GetIntForFlag(key string) (*int, error) {
	return c.GetInt(c.getIndex(key))
}

// GetStringList returns a string list
func (c *Context) GetStringList(index int) (*[]string, error) {
	str, err := c.GetString(index)
	if err != nil {
		return nil, err
	}
	arr := strings.Split(*str, ",")
	return &arr, nil
}

// GetStringListForFlag returns a string list
func (c *Context) GetStringListForFlag(key string) (*[]string, error) {
	return c.GetStringList(c.getIndex(key))
}

// GetDictionary returns a key / value dictionary
func (c *Context) GetDictionary(index int) (*map[string]string, error) {
	dic, err := c.GetString(index)
	if err != nil {
		return nil, err
	}

	dicKeyValues := map[string]string{}
	if len(*dic) > 0 {
		keyValArr := strings.Split(*dic, ",")
		if len(keyValArr)%2 != 0 {
			return nil, errors.New("attribute keys:value count does not match")
		}
		for i := 0; i < len(keyValArr); i += 2 {
			key := keyValArr[i]
			if len(key) == 0 {
				return nil, errors.New("bad key")
			}
			val := keyValArr[i+1]
			if len(val) == 0 {
				return nil, errors.New("bad val")
			}
			dicKeyValues[key] = val
		}
	}
	return &dicKeyValues, nil
}

// GetDictionaryForFlag returns a key / value dictionary
func (c *Context) GetDictionaryForFlag(key string) (*map[string]string, error) {
	return c.GetDictionary(c.getIndex(key))
}

// Count returns the number of args
func (c *Context) Count() int {
	return len(c.Args)
}

// - private

func (c *Context) validate(argName string, obj interface{}, kind reflect.Kind) error {
	if obj == nil {
		return errors.New(argName + " is missing")
	}
	if reflect.TypeOf(obj).Kind() != reflect.Ptr {
		return errors.New(argName + " should be a pointer")
	}
	if reflect.TypeOf(obj).Elem().Kind() != kind {
		return errors.New(argName + " value has invalid format")
	}
	return nil
}

func (c *Context) getIndex(key string) int {
	for idx, arg := range c.Args {
		if strings.Compare(arg.Flag, key) == 0 {
			return idx
		}
	}
	panic("can't find value for: " + key)
}

func (c *Context) get(index int) interface{} {
	if index >= len(c.Args) {
		panic("index out of bounds")
	}
	return c.Args[index].Value
}
