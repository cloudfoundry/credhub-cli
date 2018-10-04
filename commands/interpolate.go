package commands

import (
	"fmt"
	"io/ioutil"
	"path"

	"code.cloudfoundry.org/credhub-cli/errors"
	"github.com/cloudfoundry/bosh-cli/director/template"
)

type InterpolateCommand struct {
	File   string `short:"f" long:"file"   description:"Path to the file to interpolate"`
	Prefix string `short:"p" long:"prefix" description:"Prefix to be applied to credential paths. Will not be applied to paths that start with '/'"`
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
		return errors.NewEmptyTemplateError(c.File)
	}

	initialTemplate := template.NewTemplate(fileContents)

	varGetter := variableGetter{
		clientCommand: c.ClientCommand,
		prefix:        c.Prefix,
	}

	renderedTemplate, err := initialTemplate.Evaluate(varGetter, nil, template.EvaluateOpts{ExpectAllKeys: true})
	if err != nil {
		return err
	}

	fmt.Println(string(renderedTemplate))
	return nil
}

type variableGetter struct {
	clientCommand ClientCommand
	prefix        string
}

func (v variableGetter) Get(varDef template.VariableDefinition) (interface{}, bool, error) {
	varName := varDef.Name
	if !path.IsAbs(varDef.Name) {
		varName = path.Join(v.prefix, varName)
	}

	variable, err := v.clientCommand.client.GetLatestVersion(varName)
	var result = variable.Value
	if mapString, ok := variable.Value.(map[string]interface{}); ok {
		mapInterface := map[interface{}]interface{}{}
		for k, v := range mapString {
			mapInterface[k] = v
		}
		result = mapInterface
	}

	return result, true, err
}

func (v variableGetter) List() ([]template.VariableDefinition, error) {
	// not implemented
	return []template.VariableDefinition{}, nil
}
