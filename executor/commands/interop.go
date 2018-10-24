package commands

import (
	"github.com/Microsoft/kunlun/common/errors"
	"github.com/Microsoft/kunlun/common/storage"
)

type Interop struct {
	stateStore storage.Store
}

func NewInterop(
	stateStore storage.Store,
) Interop {
	return Interop{
		stateStore: stateStore,
	}
}

func (p Interop) CheckFastFails(args []string, state storage.State) error {
	return &errors.NotImplementedError{}
}

func (p Interop) Execute(args []string, state storage.State) error {
	return &errors.NotImplementedError{}
}
