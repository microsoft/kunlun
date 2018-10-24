package commands

import (
	"github.com/Microsoft/kunlun/common/errors"
	"github.com/Microsoft/kunlun/common/storage"
)

type Promote struct {
	stateStore storage.Store
}

func NewPromote(
	stateStore storage.Store,
) Promote {
	return Promote{
		stateStore: stateStore,
	}
}

func (p Promote) CheckFastFails(args []string, state storage.State) error {
	return &errors.NotImplementedError{}
}

func (p Promote) Execute(args []string, state storage.State) error {
	return &errors.NotImplementedError{}
}
