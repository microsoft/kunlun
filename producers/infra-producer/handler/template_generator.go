package handler

import (
	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/storage"
)

type TemplateGenerator interface {
	GenerateTemplate(artifacts.Manifest, storage.State) string
}
