package kevin

import (
	"errors"
	"fmt"
	"github.com/realfabecker/kevin/internal/core/domain"
	"github.com/realfabecker/kevin/internal/core/ports"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type YmlCommandRepository struct {
	ports.Logger
}

func NewYmlCommandRepository(logger ports.Logger) ports.CommandRepository {
	return &YmlCommandRepository{logger}
}

// Get return a command by its name
func (m YmlCommandRepository) Get(name string) (*domain.Cmd, error) {
	p, err := m.List()
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}
	for _, v := range p {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("%s is not a valid command", name)
}

// List return a list of repositories
func (m YmlCommandRepository) List() ([]domain.Cmd, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var cmds = make([]domain.Cmd, 0)
	fp := filepath.Join(wd, "kevin.yml")
	if cm, err := m.source(fp); err != nil {
		return nil, fmt.Errorf("source(1): %w", err)
	} else if cm != nil {
		cmds = append(cmds, cm...)
	}

	ud, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	up := filepath.Join(ud, ".kevin", "kevin.yml")
	if cm, err := m.source(up); err != nil {
		return nil, fmt.Errorf("source(1): %w", err)
	} else if cm != nil {
		cmds = append(cmds, cm...)
	}
	return cmds, nil
}

func (m YmlCommandRepository) source(rf string) ([]domain.Cmd, error) {
	m.Logger.Debug("reading from working ref: " + rf)
	cd, err := os.ReadFile(rf)

	if err != nil && errors.Is(err, os.ErrNotExist) == false {
		return nil, err
	} else if cd == nil {
		return nil, nil
	}

	var src struct {
		Commands []domain.Cmd `yaml:"commands"`
	}
	if err := yaml.Unmarshal(cd, &src); err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	for i, v := range src.Commands {
		if v.Ref == "" {
			continue
		}

		rp := v.Ref
		if !filepath.IsAbs(v.Ref) {
			rp = filepath.Join(filepath.Dir(rf), rp)
		}

		if src.Commands[i].Pipe, err = m.source(rp); err != nil {
			return nil, err
		}
	}
	return src.Commands, nil
}
