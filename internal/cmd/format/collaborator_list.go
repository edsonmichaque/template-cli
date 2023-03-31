package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type CollaboratorList dnsimple.CollaboratorsResponse

func (a CollaboratorList) FormatJSON(opts *Opts) (io.Reader, error) {
	return formatJSON(a, opts)
}

func (a CollaboratorList) FormatYAML(opts *Opts) (io.Reader, error) {
	return formatYAML(a, opts)
}

func (a CollaboratorList) FormatTable(_ *Opts) (io.Reader, error) {
	return formatTable(a)
}

func (a CollaboratorList) formatJSON(opts *Opts) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a CollaboratorList) formatHeader() []string {
	return []string{
		"ID",
		"DOMAIN ID",
		"DOMAIN NAME",
		"USER ID",
		"USER EMAIL",
		"INVITATION",
		"CREATED AT",
		"UPDATED AT",
		"ACCEPTED AT",
	}
}

func (a CollaboratorList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(a.Data))

	domains := a.Data

	for i := range domains {
		data = append(data, map[string]string{
			"ID":          fmt.Sprintf("%d", domains[i].ID),
			"DOMAIN ID":   fmt.Sprintf("%d", domains[i].DomainID),
			"DOMAIN NAME": domains[i].DomainName,
			"USER ID":     fmt.Sprintf("%d", domains[i].UserID),
			"USER EMAIL":  domains[i].UserEmail,
			"INVITATION":  fmt.Sprintf("%t", domains[i].Invitation),
			"CREATED AT":  domains[i].CreatedAt,
			"UPDATED AT":  domains[i].UpdatedAt,
			"ACCEPTED AT": domains[i].AcceptedAt,
		})
	}

	return data
}
