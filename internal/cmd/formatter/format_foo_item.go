package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type DSRItem Foo

func (d DSRItem) FormatText(opts *Opts) (io.Reader, error) {
	keys := []string{
		"id",
		"name",
		"age",
	}

	const txtLen = 8

	values := map[string]interface{}{
		"id":   d.ID,
		"name": d.Name,
		"age":  d.Age,
	}

	titles := map[string]string{
		"id":   "ID",
		"name": "Domain ID",
		"age":  "Age",
	}

	buf := new(bytes.Buffer)
	for _, v := range keys {
		buf.WriteString(fmt.Sprintf("%-20s%v\n", titles[v]+":", values[v]))
	}

	return buf, nil
}

func (d DSRItem) FormatJSON(opts *Opts) (io.Reader, error) {
	return formatJSON(d, opts)
}

func (d DSRItem) FormatYAML(opts *Opts) (io.Reader, error) {
	return formatYAML(d, opts)
}

func (d DSRItem) formatJSON(opts *Opts) ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}
