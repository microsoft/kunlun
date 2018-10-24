package handler

import (
	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/storage"
)

type Manager interface {
	Setup(manifest artifacts.Manifest, kunlunState storage.State) error
	Apply(kunlunState storage.State) (storage.State, error)
	GetOutputs() (Outputs, error)
}
