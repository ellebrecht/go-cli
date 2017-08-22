package output

import (
	"fmt"
	"io"
	"sync"
)

// Output encapsulates terminal output
type Output struct {
	IsEnabled bool
}

var instance *Output
var once sync.Once

type ActionHandler func()

// @see http://marcio.io/2015/07/singleton-pattern-in-go/
func shared() *Output {
	once.Do(func() {
		instance = &Output{}
		instance.IsEnabled = true
	})
	return instance
}

// - Public

// Println copy of fmt.Println()
func Println(a ...interface{}) {
	if !shared().IsEnabled {
		return
	}
	fmt.Println(a...)
}

// Printf copy of fmt.Printf()
func Printf(format string, a ...interface{}) {
	if !shared().IsEnabled {
		return
	}
	fmt.Printf(format, a...)
}

// Fprintf copy of fmt.Fprintf()
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	if !shared().IsEnabled {
		return 0, nil
	}
	return fmt.Fprintf(w, format, a...)
}

// Fprint copy of fmt.Fprint()
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	if !shared().IsEnabled {
		return 0, nil
	}
	return fmt.Fprint(w, a...)
}

// SetIsEnabled true allows output to terminal
func SetIsEnabled(isEnabled bool) {
	shared().IsEnabled = isEnabled
}

// DisableForAction temporarily disables output for given action
func DisableForAction(action ActionHandler) {
	shared().IsEnabled = false
	action()
	shared().IsEnabled = true
}
