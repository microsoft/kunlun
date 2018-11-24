package terraform

import (
	"bytes"
	"errors"
	"fmt"

	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	"github.com/Microsoft/kunlun/producers/infra-producer/handler"
	"github.com/coreos/go-semver/semver"
)

type Manager struct {
	executor              Executor
	templateGenerator     TemplateGenerator
	inputGenerator        InputGenerator
	terraformOutputBuffer *bytes.Buffer
	ui                    *ui.UI
}

func NewManager(executor Executor, templateGenerator TemplateGenerator, inputGenerator InputGenerator, terraformOutputBuffer *bytes.Buffer, ui *ui.UI) Manager {
	return Manager{
		executor:              executor,
		templateGenerator:     templateGenerator,
		inputGenerator:        inputGenerator,
		terraformOutputBuffer: terraformOutputBuffer,
		ui:                    ui,
	}
}

func (m Manager) Version() (string, error) {
	return m.executor.Version()
}

func (m Manager) ValidateVersion() error {
	version, err := m.executor.Version()
	if err != nil {
		return err
	}

	currentVersion, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	minimumVersion, err := semver.NewVersion("0.11.0")
	if err != nil {
		return err
	}

	if currentVersion.LessThan(*minimumVersion) {
		return errors.New("Terraform version must be at least v0.11.0")
	}

	return nil
}

func (m Manager) Setup(manifest artifacts.Manifest, kunlunState storage.State) error {
	m.ui.Step("generating terraform template")
	template, err := m.templateGenerator.GenerateTemplate(manifest, kunlunState)
	if err != nil {
		return fmt.Errorf("Template generator generate: %s", err)
	}

	m.ui.Step("generating terraform variables")
	input, err := m.inputGenerator.GenerateInput(manifest, kunlunState)
	if err != nil {
		return fmt.Errorf("Input generator generate: %s", err)
	}

	if err := m.executor.Setup(template, input); err != nil {
		return fmt.Errorf("Executor setup: %s", err)
	}

	return m.Init(kunlunState)
}

func (m Manager) Init(kunlunState storage.State) error {
	m.ui.Step("terraform init")
	if err := m.executor.Init(); err != nil {
		return fmt.Errorf("Executor init: %s", err)
	}
	return nil
}

func (m Manager) Apply(kunlunState storage.State) (storage.State, error) {
	m.ui.Step("terraform init")
	if err := m.executor.Init(); err != nil {
		return kunlunState, fmt.Errorf("Executor init: %s", err)
	}

	m.ui.Step("terraform apply")
	err := m.executor.Apply(m.inputGenerator.Credentials(kunlunState))

	kunlunState.LatestTFOutput = readAndReset(m.terraformOutputBuffer)

	if err != nil {
		return kunlunState, fmt.Errorf("Executor apply: %s", err)
	}

	return kunlunState, nil
}

func (m Manager) Destroy(kunlunState storage.State) (storage.State, error) {
	m.ui.Step("terraform destroy")
	err := m.executor.Destroy(m.inputGenerator.Credentials(kunlunState))

	kunlunState.LatestTFOutput = readAndReset(m.terraformOutputBuffer)

	if err != nil {
		return kunlunState, fmt.Errorf("Executor destroy: %s", err)
	}

	m.ui.Step("finished destroying infrastructure")
	return kunlunState, nil
}

func (m Manager) Validate(kunlunState storage.State) (storage.State, error) {
	m.ui.Step("terraform validate")
	err := m.executor.Validate(m.inputGenerator.Credentials(kunlunState))

	kunlunState.LatestTFOutput = readAndReset(m.terraformOutputBuffer)

	if err != nil {
		return kunlunState, fmt.Errorf("Executor validate: %s", err)
	}

	return kunlunState, nil
}

func (m Manager) GetOutputs() (handler.Outputs, error) {
	tfOutputs, err := m.executor.Outputs()
	if err != nil {
		return Outputs{}, err
	}

	return Outputs{Map: tfOutputs}, nil
}

func (m Manager) IsPaved() (bool, error) {
	return m.executor.IsPaved()
}

func readAndReset(buf *bytes.Buffer) string {
	contents := buf.Bytes()
	buf.Reset()

	return string(contents)
}
