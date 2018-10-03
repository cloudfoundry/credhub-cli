package commands

import (
	"fmt"
	"io/ioutil"

	"code.cloudfoundry.org/credhub-cli/errors"
	"github.com/cloudfoundry/bosh-cli/director/template"
)

type InterpolateCommand struct {
	File string `short:"f" long:"file" description:"Path to the file to interpolate"`
	ClientCommand
}

func (c *InterpolateCommand) Execute([]string) error {
	var ()

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
}

func (v variableGetter) Get(varDef template.VariableDefinition) (interface{}, bool, error) {
	variable, err := v.clientCommand.client.GetLatestVersion(varDef.Name)
	return variable.Value, true, err
}

func (v variableGetter) List() ([]template.VariableDefinition, error) {
	// not implemented
	return []template.VariableDefinition{}, nil
}
