package apis

import yaml "gopkg.in/yaml.v2"

type Role struct {
	Name       string        `yaml:"name"`
	Vars       yaml.MapSlice `yaml:"vars"`
	BecomeUser string        `yaml:"become_user"`
}
