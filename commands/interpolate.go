package commands

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/errors"
	"github.com/cloudfoundry/bosh-cli/director/template"
)

type InterpolateCommand struct {
	File              string `short:"f" long:"file"   description:"Path to the file to interpolate"`
	Prefix            string `short:"p" long:"prefix" description:"Prefix to be applied to credential paths. Will not be applied to paths that start with '/'"`
	SkipMissingParams bool   `short:"s" long:"skip-missing" description:"allow skipping missing params"`
	ClientCommand
}

func (c *InterpolateCommand) Execute([]string) error {
	if c.File == "" {
		return errors.NewMissingInterpolateParametersError()
	}

	fileContents, err := ioutil.ReadFile(c.File)
	if err != nil {
		return err
	}

	if len(fileContents) == 0 {
		return nil
	}

	initialTemplate := template.NewTemplate(fileContents)

	credGetter := credentialGetter{
		clientCommand: c.ClientCommand,
		prefix:        c.Prefix,
		creds:         make(map[string]credentials.Credential),
	}

	results, err := c.ClientCommand.client.FindByPath("")
	if err != nil {
		return err
	}
	paths := []string{}
	for _, result := range results.Credentials {
		paths = append(paths, result.Name)
	}
	credGetter.paths = paths
	renderedTemplate, err := initialTemplate.Evaluate(credGetter, nil, template.EvaluateOpts{ExpectAllKeys: !c.SkipMissingParams})
	if err != nil {
		return err
	}

	fmt.Println(string(renderedTemplate))
	return nil
}

type credentialGetter struct {
	clientCommand ClientCommand
	prefix        string
	paths         []string
	creds         map[string]credentials.Credential
}

func (v credentialGetter) Get(varDef template.VariableDefinition) (interface{}, bool, error) {
	credName := varDef.Name
	if !path.IsAbs(varDef.Name) {
		credName = path.Join("/", v.prefix, credName)
	}
	if v.HasCredential(credName) {
		var credential credentials.Credential
		if val, ok := v.creds[credName]; ok {
			credential = val
		} else {
			val, err := v.clientCommand.client.GetLatestVersion(credName)
			if err != nil {
				return nil, false, err
			}
			credential = val
			v.creds[credName] = credential
		}
		var result = credential.Value
		if mapString, ok := credential.Value.(map[string]interface{}); ok {
			mapInterface := map[interface{}]interface{}{}
			for k, v := range mapString {
				mapInterface[k] = v
			}
			result = mapInterface
		}
		return result, true, nil
	} else {
		return nil, false, nil
	}
}

func (v credentialGetter) HasCredential(credName string) bool {
	for _, cred := range v.paths {
		if strings.EqualFold(cred, credName) {
			return true
		}
	}
	return false
}

func (v credentialGetter) List() ([]template.VariableDefinition, error) {
	// not implemented
	return []template.VariableDefinition{}, nil
}
