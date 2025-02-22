package pythonx

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed template.tpl
var tplcontent string
var tpl = template.Must(template.New("Dockerfile").Parse(tplcontent))
var DefaultPythonVersion = "python:3.13-slim"

type DockerfileVars struct {
	Version string
	Ports   []string
	// Add any future container configuration here
}

func NewDockerfileVars() *DockerfileVars {
	return &DockerfileVars{
		Version: DefaultPythonVersion,
		Ports:   []string{},
	}
}

func Dockerfile(vars *DockerfileVars) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)
	if err := tpl.Execute(buff, vars); err != nil {
		return nil, err
	}

	return buff, nil
}
