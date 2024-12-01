package ports

import "github.com/realfabecker/kevin/internal/core/domain"

type Logger interface {
	Info(message string)
	Infof(format string, a ...interface{})
	Error(message string)
	Errorf(format string, a ...interface{})
	Debug(message string)
	Warn(message string)
	Fataln(v ...any)
}

type Cli interface {
	Run(cmd domain.Cmd) error
}

type CommandRepository interface {
	Get(name string) (*domain.Cmd, error)
	List() ([]domain.Cmd, error)
}

type ParallelRunner interface {
	Run(command string, pll int, mFlags []map[string]string)
}

type FlagListReader interface {
	Read(p string) ([]map[string]string, error)
}

type ScriptRender interface {
	Render(cmd *domain.Cmd) (string, error)
}
