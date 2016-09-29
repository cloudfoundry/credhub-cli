package commands

import (
	. "github.com/pivotal-cf/credhub-cli/models"
	"github.com/pivotal-cf/credhub-cli/util"
	"strings"
	"encoding/json"
)

type PrinterFactory interface {
	PrintableSecret() string
}

func NewPrinterFactory(s Secret) PrinterFactory {
	return printerFactory{
		secret: s,
		printFunc: allPrinters[s.SecretBody.ContentType],
	}
}

type printerList map[string](func(secret Secret) []string)

var allPrinters = printerList{
	"value": valuePrinter,
	"password": valuePrinter,
	"certificate": certificatePrinter,
	"ssh": sshPrinter,
}

type printerFactory struct {
	secret    Secret
	printFunc func(secret Secret) []string
}

func (p printerFactory) PrintableSecret() string {
	lines := []string{}

	secretBody := p.secret.SecretBody
	lines = append(lines,
		util.BuildLineOfFixedLength("Type:", secretBody.ContentType),
		util.BuildLineOfFixedLength("Name:", p.secret.Name),
	)
	lines = append(lines, p.printFunc(p.secret)...)

	lines = append(lines, util.BuildLineOfFixedLength("Updated:", secretBody.UpdatedAt))

	return strings.Join(lines, "\n")
}

func valuePrinter(secret Secret) []string {
	lines := []string{}
	value := secret.SecretBody.Value.(string)
	lines = append(lines, util.BuildLineOfFixedLength("Value:", value))
	return lines
}

func certificatePrinter(secret Secret) []string {
	// We are marshaling again here because there isn't a simple way
	//// to convert map[string]interface{} to a Certificate struct
	lines := []string{}
	json_cert, _ := json.Marshal(secret.SecretBody.Value)
	cert := Certificate{}
	json.Unmarshal(json_cert, &cert)
	lines = append(lines, cert.StringLines()...)
	return lines
}

func sshPrinter(secret Secret) []string {
	lines := []string{}
	json_ssh, _ := json.Marshal(secret.SecretBody.Value)
	ssh := Ssh{}
	json.Unmarshal(json_ssh, &ssh)
	return append(lines, ssh.StringLines()...)
}