package cmd

import (
	"github.com/realfabecker/kevin/internal/adapters/kevin"
	"github.com/realfabecker/kevin/internal/adapters/logger"
	"github.com/realfabecker/kevin/internal/cmd/pll"
	"github.com/realfabecker/kevin/internal/cmd/run"
	"github.com/realfabecker/kevin/internal/core/domain"
	"github.com/spf13/cobra"
)

var (
	console = logger.NewConsoleLogger()
)

func readConfig() ([]domain.Cmd, error) {
	repo := kevin.NewYmlCommandRepository(console)
	return repo.List()
}

func newRootCmd(cmds []domain.Cmd) *cobra.Command {
	var root = &cobra.Command{
		Use:           "kevin [command]",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	root.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})
	root.AddCommand(pll.NewRunCmd())
	if len(cmds) > 0 {
		run.AttachCmd(root, cmds)
	}
	return root
}

func Execute() {
	defer func() {
		if r := recover(); r != nil {
			console.Fataln(r)
		}
	}()

	cmds, err := readConfig()
	if err != nil {
		console.Fataln(err)
	}

	if err := newRootCmd(cmds).Execute(); err != nil {
		console.Fataln(err)
	}
}
