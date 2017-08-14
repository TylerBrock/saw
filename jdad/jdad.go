package jdad

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"

	"github.com/fatih/color"
)

type formatter struct {
	KeyColor        *color.Color
	StringColor     *color.Color
	BoolColor       *color.Color
	NumberColor     *color.Color
	NullColor       *color.Color
	StringMaxLength int
	DisabledColor   bool
	Indent          int
}

var fmter = formatter{
	KeyColor:        color.New(color.FgWhite),
	StringColor:     color.New(color.FgGreen),
	BoolColor:       color.New(color.FgYellow),
	NumberColor:     color.New(color.FgCyan),
	NullColor:       color.New(color.FgMagenta),
	StringMaxLength: 0,
	DisabledColor:   false,
	Indent:          0,
}

func serializeMap(m map[string]interface{}, buf *bytes.Buffer) {
	buf.WriteString("{ ")
	remaining := len(m)

	if remaining == 0 {
		buf.WriteString("{}")
		return
	}

	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		buf.WriteString(fmter.KeyColor.Sprintf("\"%s\"", key))
		buf.WriteString(": ")
		serializeValue(m[key], buf)
		remaining--
		if remaining != 0 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(" }")
}

func serializeArray(a []interface{}, buf *bytes.Buffer) {
	if len(a) == 0 {
		buf.WriteString("[]")
		return
	}
	buf.WriteString("[ ")
	for i, v := range a {
		serializeValue(v, buf)
		if i < len(a)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(" ]")
}

func serializeValue(val interface{}, buf *bytes.Buffer) {
	switch v := val.(type) {
	case map[string]interface{}:
		serializeMap(v, buf)
	case []interface{}:
		serializeArray(v, buf)
	case string:
		b, _ := json.Marshal(v)
		buf.WriteString(fmter.StringColor.SprintFunc()(string(b)))
	case float64:
		buf.WriteString(fmter.NumberColor.SprintFunc()(strconv.FormatFloat(v, 'f', -1, 64)))
	case bool:
		buf.WriteString(fmter.BoolColor.SprintFunc()(strconv.FormatBool(v)))
	case nil:
		buf.WriteString(fmter.NullColor.SprintFunc()("null"))
	}
}

func Serialize(jsonMap map[string]interface{}) string {
	buffer := bytes.Buffer{}
	serializeMap(jsonMap, &buffer)
	return buffer.String()
}
