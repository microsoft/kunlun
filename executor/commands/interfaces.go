package commands

import "github.com/Microsoft/kunlun/common/storage"

type logger interface {
	Step(string, ...interface{})
	Printf(string, ...interface{})
	Println(string)
	Prompt(string) bool
}

type StateStore interface {
	Set(state storage.State) error
	GetVarsDir() (string, error)
	GetArtifactsDir() (string, error)
}
