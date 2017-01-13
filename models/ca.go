package models

import (
	"strings"

	"github.com/pivotal-cf/credhub-cli/util"
)

type Ca struct {
	CaBody CaBody
}

func NewCa(name string, caBody CaBody) Ca {
	ca := Ca{
		CaBody: caBody,
	}
	ca.CaBody.Name = name
	return ca
}

func (ca Ca) Terminal() string {
	lines := []string{}

	caBody := ca.CaBody

	if caBody.Value.Certificate != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Certificate:", caBody.Value.Certificate))
	}

	if caBody.Value.PrivateKey != "" {
		lines = append(lines, util.BuildLineOfFixedLength("Private Key:", caBody.Value.PrivateKey))
	}

	return util.Header(caBody.SecretType, ca.CaBody.Name) + strings.Join(lines, "\n") + "\n" + util.Footer(ca.CaBody.VersionCreatedAt)
}

func (ca Ca) Json() string {
	return prettyPrintJson(
		map[string]interface{}{
			"type":               ca.CaBody.SecretType,
			"version_created_at": ca.CaBody.VersionCreatedAt,
			"certificate":        ca.CaBody.Value.Certificate,
			"private_key":        ca.CaBody.Value.PrivateKey,
		})
}
