package bootstrap

import "flag"

type CommandFlags struct {
	Conf string
	Env  string
}

func NewCommandFlags() *CommandFlags {
	return &CommandFlags{
		Conf: "",
		Env:  "",
	}
}

func (f *CommandFlags) Init() {
	flag.StringVar(&f.Conf, "conf", "./configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&f.Env, "env", "dev", "runtime environment, eg: -env dev")
}
