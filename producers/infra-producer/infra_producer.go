package apis

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/errors"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	"github.com/Microsoft/kunlun/producers/infra-producer/handler"
	"github.com/Microsoft/kunlun/producers/infra-producer/tfhandler/terraform"
	"github.com/spf13/afero"
)

type InfraProducer struct {
	manager handler.Manager
}

func NewInfraProducer(stateStore storage.Store, handlerType string, debug bool) (InfraProducer, error) {
	log.SetFlags(0)

	ui := ui.NewUI(os.Stdout, os.Stdin)

	fs := afero.NewOsFs()
	afs := &afero.Afero{Fs: fs}

	if handlerType == handler.TerraformHandlerType {
		terraformOutputBuffer := bytes.NewBuffer([]byte{})
		terraformDir, _ := stateStore.GetTerraformDir()
		dotTerraformDir := filepath.Join(terraformDir, ".terraform")
		bufferingCLI := terraform.NewCLI(terraformOutputBuffer, terraformOutputBuffer, dotTerraformDir)
		var (
			terraformCLI terraform.CLI
			out          io.Writer
		)
		if debug {
			errBuffer := io.MultiWriter(os.Stderr, terraformOutputBuffer)
			terraformCLI = terraform.NewCLI(errBuffer, terraformOutputBuffer, dotTerraformDir)
			out = os.Stdout
		} else {
			terraformCLI = bufferingCLI
			out = ioutil.Discard
		}
		terraformExecutor := terraform.NewExecutor(terraformCLI, bufferingCLI, stateStore, afs, debug, out)

		inputGenerator := terraform.NewInputGenerator()
		templateGenerator := terraform.NewTemplateGenerator()
		terraformManager := terraform.NewManager(terraformExecutor, templateGenerator, inputGenerator, terraformOutputBuffer, ui)

		return InfraProducer{
			manager: terraformManager,
		}, nil
	} else if handlerType == handler.ARMTemplateHandlerType {
		return InfraProducer{}, &errors.NotImplementedError{}
	} else {
		return InfraProducer{}, &errors.NotSupportedError{}
	}
}

func (ip InfraProducer) Setup(manifest artifacts.Manifest, state storage.State) error {
	return ip.manager.Setup(manifest, state)
}

func (ip InfraProducer) Apply(state storage.State) error {
	_, err := ip.manager.Apply(state)
	if err != nil {
		return err
	}

	return nil
}

func (ip InfraProducer) GetOutputs() (string, error) {
	outputs, err := ip.manager.GetOutputs()
	if err != nil {
		return "", err
	}

	contents, err := handler.ToOutputsOpsFile(outputs)
	if err != nil {
		return "", err
	}

	return contents, nil
}
