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
	d, err := m.source()
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	var src struct {
		Commands []domain.Cmd `yaml:"commands"`
	}
	if err := yaml.Unmarshal(d, &src); err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	return src.Commands, nil
}

// Return the configuration source
func (m YmlCommandRepository) source() ([]byte, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fp := filepath.Join(wd, "kevin.yml")
	m.Logger.Debug("reading from working path: " + fp)
	if cd, err := os.ReadFile(fp); err != nil && errors.Is(err, os.ErrNotExist) == false {
		return nil, err
	} else if cd != nil {
		return cd, nil
	}

	ud, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	up := filepath.Join(ud, "kevin.yml")
	m.Logger.Debug("reading from user path: " + up)
	return os.ReadFile(filepath.Join(ud, "kevin.yml"))
}
