package template

import (
	"fmt"

	"github.com/cloudfoundry/bosh-cli/director/template"

	yaml "gopkg.in/yaml.v2"
)

type MissingVariableError struct {
	name string
}

func (e MissingVariableError) Error() string {
	return fmt.Sprintf("unbound variable in template: '%s'", e.name)
}

type FileVarsSource struct {
	ParamsContent []byte
}

func (s *FileVarsSource) Evaluate(content []byte) ([]byte, error) {
	tpl := template.NewTemplate(content)

	var vars template.StaticVariables
	err := yaml.Unmarshal(s.ParamsContent, &vars)
	if err != nil {
		return nil, err
	}

	bytes, err := tpl.Evaluate(vars, nil, template.EvaluateOpts{ExpectAllKeys: true})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
