package configuration

import "github.com/Microsoft/kunlun/common/storage"

type GlobalConfiguration struct {
	StateDir string
	Debug    bool
	Name     string
}

type StringSlice []string

func (s StringSlice) ContainsAny(targets ...string) bool {
	for _, target := range targets {
		for _, element := range s {
			if element == target {
				return true
			}
		}
	}
	return false
}

type Configuration struct {
	Global               GlobalConfiguration
	Command              string
	SubcommandFlags      StringSlice
	State                storage.State
	ShowCommandHelp      bool
	CommandModifiesState bool
}
