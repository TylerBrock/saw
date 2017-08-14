package jdad

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const StartMap = '{'
const EndMap = '}'
const StartArray = '['
const EndArray = ']'

type Formatter struct {
	KeyColor        *color.Color
	StringColor     *color.Color
	BoolColor       *color.Color
	NumberColor     *color.Color
	NullColor       *color.Color
	StringMaxLength int
	DisabledColor   bool
	Indent          int
}

var fmter = Formatter{
	KeyColor:        color.New(color.FgWhite),
	StringColor:     color.New(color.FgGreen),
	BoolColor:       color.New(color.FgYellow),
	NumberColor:     color.New(color.FgCyan),
	NullColor:       color.New(color.FgMagenta),
	StringMaxLength: 0,
	DisabledColor:   false,
	Indent:          0,
}

type sprinter struct {
	Key    func(...interface{}) string
	String func(...interface{}) string
	Bool   func(...interface{}) string
	Number func(...interface{}) string
	Null   func(...interface{}) string
}

var sptr = sprinter{
	Key:    fmter.KeyColor.SprintFunc(),
	String: fmter.StringColor.SprintFunc(),
	Bool:   fmter.BoolColor.SprintFunc(),
	Number: fmter.NumberColor.SprintFunc(),
	Null:   fmter.NullColor.SprintFunc(),
}

func (f *Formatter) writeIndent(buf *bytes.Buffer, depth int) {
	buf.WriteString(strings.Repeat(" ", f.Indent*depth))
}

func serializeMap(m map[string]interface{}, buf *bytes.Buffer, depth int) {
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

	buf.WriteString("{")
	for _, key := range keys {
		buf.WriteString(fmter.KeyColor.Sprintf("\"%s\"", key))
		buf.WriteString(": ")
		serializeValue(m[key], buf, 0)
		remaining--
		if remaining != 0 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("}")
}

func serializeArray(a []interface{}, buf *bytes.Buffer, depth int) {
	if len(a) == 0 {
		buf.WriteString("[]")
		return
	}

	buf.WriteString("[")
	for i, v := range a {
		serializeValue(v, buf, 0)
		if i < len(a)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
}

func serializeValue(val interface{}, buf *bytes.Buffer, indent int) {
	switch v := val.(type) {
	case map[string]interface{}:
		serializeMap(v, buf, 0)
	case []interface{}:
		serializeArray(v, buf, 0)
	case string:
		b, _ := json.Marshal(v)
		buf.WriteString(sptr.String(string(b)))
	case float64:
		buf.WriteString(sptr.Number(strconv.FormatFloat(v, 'f', -1, 64)))
	case bool:
		buf.WriteString(sptr.Bool(strconv.FormatBool(v)))
	case nil:
		buf.WriteString(sptr.Null("null"))
	}
}

func Serialize(jsonMap map[string]interface{}) string {
	buffer := bytes.Buffer{}
	serializeMap(jsonMap, &buffer, 0)
	return buffer.String()
}
