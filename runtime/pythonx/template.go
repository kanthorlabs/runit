package pythonx

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed template.tpl
var tplcontent string
var tpl = template.Must(template.New("Dockerfile").Parse(tplcontent))

type DockerfileVars struct {
	Version string
	Ports   []string
}

var DefaultVars = &DockerfileVars{
	Version: "python:3.13-slim",
	Ports:   []string{"8080"},
}

func Dockerfile(vars *DockerfileVars) (*bytes.Buffer, error) {
	buff := new(bytes.Buffer)
	if err := tpl.Execute(buff, vars); err != nil {
		return nil, err
	}

	return buff, nil
}
