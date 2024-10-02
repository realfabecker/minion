package cmd

import (
	"github.com/realfabecker/kevin/internal/adapters/kevin"
	"github.com/realfabecker/kevin/internal/adapters/logger"
	"github.com/realfabecker/kevin/internal/cmd/pll"
	"github.com/realfabecker/kevin/internal/cmd/run"
	"github.com/realfabecker/kevin/internal/core/domain"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func readConfig() ([]domain.Cmd, error) {
	echo := logger.NewConsoleLogger("kevin", os.Stdout)
	repo := kevin.NewYmlCommandRepository(echo)
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
	root.SetHelpCommand(&cobra.Command{Hidden: true})
	root.AddCommand(pll.NewRunCmd())
	if len(cmds) > 0 {
		run.AttachCmd(root, cmds)
	}
	return root
}

func Execute() {
	cmds, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
	if err := newRootCmd(cmds).Execute(); err != nil {
		log.Fatalln(err)
	}
}
