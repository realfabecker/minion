package domain

type Cmd struct {
	Name   string `yaml:"name"`
	Parent string `yaml:"parent"`
	Short  string `yaml:"short"`
	Cmd    string `yaml:"cmd"`
	Shell  string `yaml:"shell"`
	Flags  []Flag `yaml:"flags"`
}

type Flag struct {
	Name     string `yaml:"name"`
	Short    string `yaml:"short"`
	Value    string `yaml:"value"`
	Usage    string `yaml:"usage"`
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
