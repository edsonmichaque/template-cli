package formatter

import (
	"encoding/json"
	"fmt"
	"io"
)

type Foo struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  string `json:"age"`
}

type FooList []Foo

func (a FooList) FormatJSON(opts *Opts) (io.Reader, error) {
	return formatJSON(a, opts)
}

func (a FooList) FormatYAML(opts *Opts) (io.Reader, error) {
	return formatYAML(a, opts)
}

func (a FooList) FormatTable(_ *Opts) (io.Reader, error) {
	return formatTable(a)
}

func (a FooList) formatJSON(opts *Opts) ([]byte, error) {
	return json.MarshalIndent(a, "", "  ")
}

func (a FooList) formatHeader() []string {
	return []string{
		"ID",
		"NAME",
		"AGE",
	}
}

func (a FooList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(a))

	fooList := a

	const txtLen = 10

	for i := range fooList {
		data = append(data, map[string]string{
			"ID":   fmt.Sprintf("%d", fooList[i].ID),
			"NAME": fooList[i].Name,
			"AGE":  fooList[i].Age,
		})
	}

	return data
}

func truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}

	return s[:length] + "..."
}
