package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dnsimple/dnsimple-go/dnsimple"
)

type DSRList dnsimple.DelegationSignerRecordsResponse

func (a DSRList) FormatJSON(opts *Opts) (io.Reader, error) {
	return formatJSON(a, opts)
}

func (a DSRList) FormatYAML(opts *Opts) (io.Reader, error) {
	return formatYAML(a, opts)
}

func (a DSRList) FormatTable(_ *Opts) (io.Reader, error) {
	return formatTable(a)
}

func (a DSRList) formatJSON(opts *Opts) ([]byte, error) {
	return json.MarshalIndent(a.Data, "", "  ")
}

func (a DSRList) formatHeader() []string {
	return []string{
		"ID",
		"DOMAIN ID",
		"ALGORITHM",
		"DIGEST",
		"DIGEST TYPE",
		"KEYTAG",
		"PUBLIC KEY",
		"CREATED AT",
		"UPDATED AT",
	}
}

func (a DSRList) formatRows() []map[string]string {
	data := make([]map[string]string, 0, len(a.Data))

	dsr := a.Data

	const txtLen = 10

	for i := range dsr {
		data = append(data, map[string]string{
			"ID":          fmt.Sprintf("%d", dsr[i].ID),
			"DOMAIN ID":   fmt.Sprintf("%d", dsr[i].DomainID),
			"ALGORITHM":   dsr[i].Algorithm,
			"DIGEST":      truncate(dsr[i].Digest, txtLen),
			"DIGEST TYPE": dsr[i].DigestType,
			"KEYTAG":      dsr[i].Keytag,
			"PUBLIC KEY":  truncate(dsr[i].PublicKey, txtLen),
			"CREATED AT":  dsr[i].CreatedAt,
			"UPDATED AT":  dsr[i].UpdatedAt,
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
