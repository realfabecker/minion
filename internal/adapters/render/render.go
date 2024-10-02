package render

import (
	"bytes"
	"fmt"
	"github.com/realfabecker/kevin/internal/core/domain"
	"github.com/realfabecker/kevin/internal/core/ports"
	"html/template"
)

type ScriptRender struct{}

func NewScriptRender() ports.ScriptRender {
	return &ScriptRender{}
}

func (s *ScriptRender) Render(cmd *domain.Cmd) (string, error) {
	t, err := template.New("script").Parse(cmd.Cmd)
	if err != nil {
		return "", fmt.Errorf("unable do parse template: %w ", err)
	}
	buff := bytes.Buffer{}
	if err := t.Execute(&buff, cmd); err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}
	return buff.String(), nil
}
