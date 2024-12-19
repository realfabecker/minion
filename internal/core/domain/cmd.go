package domain

import "os"

type Cmd struct {
	Name   string `yaml:"name"`
	Parent string `yaml:"parent"`
	Short  string `yaml:"short"`
	Cmd    string `yaml:"cmd"`
	Flags  []Flag `yaml:"flags"`
	Args   []Arg  `yaml:"args"`
	Lang   string `yaml:"lang"`
	Ref    string `yaml:"ref"`
	Pipe   []Cmd
}

type Flag struct {
	Name     string `yaml:"name"`
	Short    string `yaml:"short"`
	Value    string `yaml:"value"`
	Usage    string `yaml:"usage"`
	Required bool   `yaml:"required"`
}

type Arg struct {
	Name     string `yaml:"name"`
	Value    string `yaml:"value"`
	Required bool   `yaml:"required"`
}

func (c *Cmd) SetFlag(flag string, value string) {
	for i, f := range c.Flags {
		if f.Name == flag {
			c.Flags[i].Value = value
			break
		}
	}
}

func (c *Cmd) GetFlag(flag string) string {
	for _, f := range c.Flags {
		if f.Name == flag {
			return f.Value
		}
	}
	return ""
}

func (c *Cmd) SetArg(arg string, value string) {
	for i, a := range c.Args {
		if a.Name == arg {
			c.Args[i].Value = value
			break
		}
	}
}

func (c *Cmd) GetArg(arg string) string {
	for _, a := range c.Args {
		if a.Name == arg {
			return a.Value
		}
	}
	return ""
}

func (c *Cmd) GetNofRequiredArgs() int {
	total := 0
	for _, a := range c.Args {
		if a.Required {
			total += 1
		}
	}
	return total
}

func (c *Cmd) GetEnv(env string) string {
	return os.Getenv(env)
}
