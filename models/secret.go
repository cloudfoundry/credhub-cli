package models

import (
	"encoding/json"
	"strings"
)

type Secret struct {
	Name       string
	SecretBody SecretBody
}

func NewSecret(name string, secretBody SecretBody) Item {
	return Secret{
		Name:       name,
		SecretBody: secretBody,
	}
}

func (secret Secret) String() string {
	lines := []string{}

	secretBody := secret.SecretBody
	lines = append(lines,
		buildLineOfFixedLength("Type:", secretBody.ContentType),
		buildLineOfFixedLength("Name:", secret.Name),
	)

	if secretBody.ContentType == "value" || secretBody.ContentType == "password" {
		value := secretBody.Value.(string)
		lines = append(lines, buildLineOfFixedLength("Value:", value))
	} else if secretBody.ContentType == "ssh" {
		json_ssh, _ := json.Marshal(secretBody.Value)
		ssh := Ssh{}
		json.Unmarshal(json_ssh, &ssh)
		if ssh.PublicKey != "" {
			lines = append(lines, buildLineOfFixedLength("Public Key:", ssh.PublicKey))
		}
		if ssh.PrivateKey != "" {
			lines = append(lines, buildLineOfFixedLength("Private Key:", ssh.PrivateKey))
		}
	} else {
		// We are marshaling again here because there isn't a simple way
		// to convert map[string]interface{} to a Certificate struct
		json_cert, _ := json.Marshal(secretBody.Value)
		cert := Certificate{}
		json.Unmarshal(json_cert, &cert)
		if cert.Ca != "" {
			lines = append(lines, buildLineOfFixedLength("Ca:", cert.Ca))
		}

		if cert.Certificate != "" {
			lines = append(lines, buildLineOfFixedLength("Certificate:", cert.Certificate))
		}

		if cert.PrivateKey != "" {
			lines = append(lines, buildLineOfFixedLength("Private Key:", cert.PrivateKey))
		}
	}

	lines = append(lines, buildLineOfFixedLength("Updated:", secretBody.UpdatedAt))

	return strings.Join(lines, "\n")
}
