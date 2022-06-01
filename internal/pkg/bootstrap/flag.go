package bootstrap

import "flag"

type CommandFlags struct {
	Conf string
}

func NewCommandFlags() *CommandFlags {
	return &CommandFlags{
		Conf: "",
	}
}

func (f *CommandFlags) Init() {
	flag.StringVar(&f.Conf, "conf", "./configs", "config path, eg: -conf config.yaml")
}
