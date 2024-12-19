package run

import (
	"fmt"
	"github.com/realfabecker/kevin/internal/adapters/logger"
	"github.com/realfabecker/kevin/internal/adapters/render"
	"github.com/realfabecker/kevin/internal/adapters/runner"
	"github.com/realfabecker/kevin/internal/core/domain"
	"github.com/spf13/cobra"
)

var DryRun bool

func newSubCmd(c domain.Cmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.Name,
		Short: c.Short,
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, f := range c.Flags {
				v, _ := cmd.Flags().GetString(f.Name)
				c.SetFlag(f.Name, v)
			}
			if len(args) > 0 && len(args) == len(c.Args) {
				for i, a := range args {
					c.Args[i].Value = a
				}
			}

			if len(args) != c.GetNofRequiredArgs() {
				return fmt.Errorf("you must supply at least %d arguments for this command", c.GetNofRequiredArgs())
			}

			rn := runner.New(runner.NewCliOpts{
				Logger: logger.NewConsoleLogger(),
				Render: render.NewScriptRender(),
			})
			return rn.Run(&c, DryRun)
		},
	}
	for _, f := range c.Flags {
		if f.Short == "" {
			f.Short = f.Name[:1]
		}
		cmd.Flags().StringP(f.Name, f.Short, f.Value, f.Usage)
		if f.Required {
			_ = cmd.MarkFlagRequired(f.Name)
		}
	}
	return cmd
}

func AttachCmd(root *cobra.Command, cmds []domain.Cmd) {
	var m = make(map[string]*cobra.Command)
	for _, v := range cmds {
		func(c domain.Cmd) {
			if c.Pipe != nil {
				groupCmd := &cobra.Command{
					Use:   c.Name,
					Short: c.Short,
				}
				AttachCmd(groupCmd, c.Pipe)
				root.AddCommand(groupCmd)
			} else {
				xmd := newSubCmd(c)
				if c.Parent != "" {
					if _, ok := m[c.Parent]; ok {
						m[c.Parent].RunE = nil
						m[c.Parent].AddCommand(xmd)
					} else {
						m[c.Parent] = &cobra.Command{
							Use: c.Parent,
						}
						m[c.Parent].AddCommand(xmd)
						root.AddCommand(m[c.Parent])
					}
				} else {
					root.AddCommand(xmd)
				}
				m[c.Name] = xmd
			}

		}(v)
	}
	root.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "run in dry run mode")
}
