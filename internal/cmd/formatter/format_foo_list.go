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

func (f FooList) FormatJSON(opts *Opts) (io.Reader, error) {
	return formatJSON(f, opts)
}

func (f FooList) FormatYAML(opts *Opts) (io.Reader, error) {
	return formatYAML(f, opts)
}

func (f FooList) FormatTable(_ *Opts) (io.Reader, error) {
	return formatTable(f)
}

func (f FooList) formatJSON(opts *Opts) ([]byte, error) {
	return json.MarshalIndent(f, "", "  ")
}

func (f FooList) formatHeader() []string {
	return []string{
		"ID",
		"NAME",
		"AGE",
	}
}

func (f FooList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(f))

	fooList := f

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
