package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type bootstrapLogger interface {
	Println(message string)
}

type StateBootstrap struct {
	bootstrapLogger bootstrapLogger
	klVersion       string
}

func NewStateBootstrap(bootstrapLogger bootstrapLogger, klVersion string) StateBootstrap {
	return StateBootstrap{
		bootstrapLogger: bootstrapLogger,
		klVersion:       klVersion,
	}
}

func (b StateBootstrap) GetState(dir string) (State, error) {
	_, err := os.Stat(dir)
	if err != nil {
		return State{}, err
	}

	file, err := os.Open(filepath.Join(dir, STATE_FILE))
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return State{}, err
	}

	state := State{}
	err = json.NewDecoder(file).Decode(&state)
	if err != nil {
		return state, err
	}
	return state, nil
}
