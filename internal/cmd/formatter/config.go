package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/edsonmichaque/template-cli/internal/config"
)

type Config struct {
	Name  string       `json:"name"`
	Type  reflect.Type `json:"type"`
	Value string       `json:"value"`
}

func (c Config) MarshalJSON() ([]byte, error) {
	config := struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	}{
		Name:  c.Name,
		Type:  fmt.Sprintf("%v", c.Type),
		Value: c.Value,
	}

	return json.Marshal(config)
}

type ConfigList []Config

func ToConfigList(c *config.Config) ConfigList {
	return []Config{
		{
			Name:  "account",
			Type:  reflect.TypeOf(c.Account),
			Value: c.Account,
		},
		{
			Name:  "access-token",
			Type:  reflect.TypeOf(c.AccessToken),
			Value: c.AccessToken,
		},
		{
			Name:  "base-url",
			Type:  reflect.TypeOf(c.BaseURL),
			Value: c.BaseURL,
		},
	}
}

func (f ConfigList) FormatJSON(opts *Opts) (io.Reader, error) {
	return formatJSON(f, opts)
}

func (f ConfigList) FormatYAML(opts *Opts) (io.Reader, error) {
	return formatYAML(f, opts)
}

func (f ConfigList) FormatTable(_ *Opts) (io.Reader, error) {
	return formatTable(f)
}

func (f ConfigList) formatJSON(opts *Opts) ([]byte, error) {
	return json.MarshalIndent(f, "", "  ")
}

func (f ConfigList) formatHeader() []string {
	return []string{
		"NAME",
		"TYPE",
		"VALUE",
	}
}

func (f ConfigList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(f))

	fooList := f

	const txtLen = 10

	for i := range fooList {
		data = append(data, map[string]string{
			"NAME":  fooList[i].Name,
			"TYPE":  fmt.Sprintf("%v", fooList[i].Type),
			"VALUE": fooList[i].Value,
		})
	}

	return data
}
