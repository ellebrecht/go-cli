package cli

import (
	"fmt"
	"strconv"
	"strings"
)

// Option represents a cli option (flag)
type Option struct {
	Name           string
	Type           OptionType
	Description    string
	Flag           string
	Aliases        []string
	IsSecure       bool
	DefaultValue   interface{}
	Value          interface{} // parsed value
	possibleValues []interface{}
	Extension      interface{}
}

// OptionType represents the type of an option
type OptionType int

const (
	// OptionTypeString represents a string flag: -x="string"
	OptionTypeString OptionType = iota

	// OptionTypeBool represents a bool flag: -x=true.
	OptionTypeBool

	// OptionTypeInt represents an int flag: -x=1.
	OptionTypeInt
)

type OptionsAlpha []*Option

func (a OptionsAlpha) Len() int           { return len(a) }
func (a OptionsAlpha) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OptionsAlpha) Less(i, j int) bool { return a[i].Name < a[j].Name }

// StringValue returns string value of option
func (o *Option) StringValue() string {
	if o.Value == nil {
		return "<nil>"
	}
	switch o.Type {
	case OptionTypeBool:
		return strconv.FormatBool(*(o.Value.(*bool)))
	case OptionTypeInt:
		return strconv.Itoa(*(o.Value.(*int)))
	case OptionTypeString:
		return *(o.Value.(*string))
	default:
		return fmt.Sprintf("%v", o.Value)
	}
}

// - private

func (o *Option) allFlags() []string {
	return append(o.Aliases, o.Flag)
}

func (o *Option) defaultValue() interface{} {
	if o.DefaultValue != nil {
		return o.DefaultValue
	}
	switch o.Type {
	case OptionTypeString:
		var a string
		return a
	case OptionTypeBool:
		var a bool
		return a
	case OptionTypeInt:
		var a int
		return a
	default:
		panic("unhandled switch statement")
	}
}

func (o *Option) valueCandidates() []interface{} {
	candidates := []interface{}{}
	for _, possibleVal := range o.possibleValues {
		switch o.Type {
		case OptionTypeBool:
			val := possibleVal.(*bool)
			if *val != o.defaultValue().(bool) {
				candidates = append(candidates, val)
			}
		case OptionTypeInt:
			val := possibleVal.(*int)
			if *val != o.defaultValue().(int) {
				candidates = append(candidates, val)
			}
		case OptionTypeString:
			val := possibleVal.(*string)
			if strings.Compare(*val, o.defaultValue().(string)) != 0 {
				candidates = append(candidates, val)
			}
		default:
			panic("unhandled switch statement")
		}
	}
	return candidates
}
