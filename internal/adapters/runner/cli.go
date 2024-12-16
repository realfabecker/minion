package runner

import (
	"bytes"
	"fmt"
	"github.com/realfabecker/kevin/internal/core/domain"
	"github.com/realfabecker/kevin/internal/core/ports"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Cli struct {
	Logger ports.Logger
	Render ports.ScriptRender
}

type NewCliOpts struct {
	Logger ports.Logger
	Render ports.ScriptRender
}

type CliRunOpts struct {
	Command   string
	Arguments []string
	Workdir   string
	Env       []string
	LogFile   string
	Attach    bool
}

func New(opts NewCliOpts) *Cli {
	return &Cli{
		Logger: opts.Logger,
		Render: opts.Render,
	}
}

func (c *Cli) Run(cmd *domain.Cmd, dryRun bool) error {
	script, err := c.Render.Render(cmd)
	if err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}

	if dryRun {
		fmt.Print(script)
		return nil
	}

	if cmd.Lang != "" {
		f, err := os.CreateTemp(os.TempDir(), "script-")
		if err != nil {
			return fmt.Errorf("unable to write temp script: %w", err)
		}
		defer func() {
			if err := os.Remove(f.Name()); err != nil {
				c.Logger.Error(err.Error())
			}
		}()
		if _, err := f.Write([]byte(script)); err != nil {
			return fmt.Errorf("uanble to write temp script: %w", err)
		}
		return c.runE(CliRunOpts{
			Command:   cmd.Lang,
			Attach:    true,
			Arguments: []string{f.Name()},
		})
	}

	if runtime.GOOS == "windows" {
		return c.runE(CliRunOpts{
			Command:   "cmd",
			Attach:    true,
			Arguments: []string{"/c", strings.TrimSpace(script)},
		})
	}
	if runtime.GOOS == "linux" {
		return c.runE(CliRunOpts{
			Command:   "/bin/bash",
			Attach:    true,
			Arguments: []string{"-c", strings.TrimSpace(script)},
		})
	}
	return fmt.Errorf("unsupported runtime: %s", runtime.GOOS)

}

func (c *Cli) runB(opts CliRunOpts) ([]byte, error) {
	cmd := exec.Command(opts.Command, opts.Arguments...)
	cmd.Dir = opts.Workdir
	cmd.Env = opts.Env
	var outb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = os.Stderr
	if c.Logger != nil {
		c.Logger.Debug(fmt.Sprintf("%s\n", cmd.String()))
	}
	err := cmd.Run()
	return outb.Bytes(), err
}

func (c *Cli) runE(opts CliRunOpts) error {
	cmd := exec.Command(opts.Command, opts.Arguments...)
	cmd.Dir = opts.Workdir
	cmd.Env = opts.Env
	if opts.Attach {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if c.Logger != nil {
		c.Logger.Debug(fmt.Sprintf("%s\n", cmd.String()))
	}
	return cmd.Wait()
}
