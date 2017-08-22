package util

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"

	log "geeny/log"
	"fmt"
)

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Pad(s string, l int) string {
	pad := l - len(s)
	if pad <= 0 {
		return s
	}
	var buf bytes.Buffer
	buf.WriteString(s)
	buf.WriteString(strings.Repeat(" ", pad))
	return buf.String()
}

func AppendLine(buf *bytes.Buffer, tableWidth []int, pos int) {
	border := [][]rune{
		{'┌', '┬', '┐'},
		{'├', '┼', '┤'},
		{'└', '┴', '┘'},
	}

	buf.WriteString("\x1b[39;1m")
	buf.WriteRune(border[pos][0])
	for i, width := range tableWidth {
		if i != 0 {
			buf.WriteRune(border[pos][1])
		}
		buf.WriteString(strings.Repeat("─", width+2))
	}
	buf.WriteRune(border[pos][2])
	buf.WriteString("\x1b[0m\n")
}

func CreateTempFile() *os.File {
	var f, err = ioutil.TempFile(os.TempDir(), "geeny-")
	if err != nil {
		panic("Unable to create temporary file " + f.Name())
	}
	log.Trace("Created temp file", f.Name())
	return f
}

func RemoveFile(f *os.File) {
	err := os.Remove(f.Name())
	if err != nil {
		panic(err)
	}
}

func TakeSliceArg(arg interface{}) []interface{} {
	slice := takeArg(arg, reflect.Slice)
	c := slice.Len()
	out := make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out
}

// - private

func takeArg(arg interface{}, kind reflect.Kind) reflect.Value {
	val := reflect.ValueOf(arg)
	if val.Kind() != kind {
		panic("Argument is not a " + kind.String())
	}
	return val
}

func StringSlice(arg interface{}) (slice []string, ok bool) {
	s, ok := arg.([]string)
	if ok {
		return s, true
	}
	l, ok := arg.([]interface{})
	if ok {
		r := make([]string, len(l))
		for i, a := range l {
			r[i] = fmt.Sprintf("%v", a)
		}
		return r, true
	}
	return nil, false
}

func StringMap(arg interface{}) (m map[string]string, ok bool) {
	s, ok := arg.(map[string]string)
	if ok {
		return s, true
	}
	l, ok := arg.(map[string]interface{})
	if ok {
		r := make(map[string]string, len(l))
		for k, v := range l {
			r[k] = fmt.Sprintf("%v", v)
		}
		return r, true
	}
	return nil, false
}

func Contains(slice []string, s string) bool {
	for _, e := range slice {
		if s == e {
			return true
		}
	}
	return false
}
