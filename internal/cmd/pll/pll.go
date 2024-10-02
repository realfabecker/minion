package pll

import (
	"github.com/realfabecker/kevin/internal/adapters/flagreader"
	"github.com/realfabecker/kevin/internal/adapters/runner"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	var f = struct {
		Workers int
		File    string
		Command string
	}{}

	cmd := &cobra.Command{
		Use: "run",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags = make([]map[string]string, 0)
			if f.File != "" {
				reader := flagreader.NewCsvFlagReader()
				if mFlags, err := reader.Read(f.File); err != nil {
					return err
				} else if mFlags != nil {
					flags = mFlags
				}
			}
			runner.NewMulti().Run(f.Command, f.Workers, flags)
			return nil
		},
	}
	cmd.Flags().IntVarP(&f.Workers, "w", "w", 1, "number of concurrent workers")
	cmd.Flags().StringVarP(&f.Command, "c", "c", "", "command to be executed")
	cmd.Flags().StringVarP(&f.File, "f", "f", "", "path to csv file with flags")
	_ = cmd.MarkFlagRequired("c")
	return cmd
}
