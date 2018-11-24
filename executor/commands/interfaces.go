package commands

import "github.com/Microsoft/kunlun/common/storage"

type StateStore interface {
	Set(state storage.State) error
	GetVarsDir() (string, error)
	GetArtifactsDir() (string, error)
}
