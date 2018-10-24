package commands

import (
	"fmt"
	"os"

	builtinmanifests "github.com/Microsoft/kunlun/artifacts/builtinmanifests"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/flags"
	"github.com/Microsoft/kunlun/common/helpers"
	"github.com/Microsoft/kunlun/common/storage"
	digesterApis "github.com/Microsoft/kunlun/digester"
	digesterCommon "github.com/Microsoft/kunlun/digester/common"
)

type Digest struct {
	stateStore   storage.Store
	envIDManager helpers.EnvIDManager

	fs fileio.Fs
}

type DiegestConfig struct {
	Name string
}

func NewDigest(
	stateStore storage.Store,
	envIDManager helpers.EnvIDManager,
	fs fileio.Fs,
) Digest {
	return Digest{
		stateStore:   stateStore,
		envIDManager: envIDManager,
		fs:           fs,
	}
}

func (p Digest) CheckFastFails(args []string, state storage.State) error {
	config, err := p.ParseArgs(args, state)
	if err != nil {
		return err
	}
	if state.EnvID != "" && config.Name != "" && config.Name != state.EnvID {
		return fmt.Errorf("The env name cannot be changed for an existing environment. Current name is %s", state.EnvID)
	}
	return nil
}

func (p Digest) ParseArgs(args []string, state storage.State) (DiegestConfig, error) {
	var (
		config DiegestConfig
	)

	digestFlags := flags.New("analyze")
	digestFlags.String(&config.Name, "name", os.Getenv("KL_ENV_NAME"))

	err := digestFlags.Parse(args)
	if err != nil {
		return DiegestConfig{}, err
	}
	return config, nil
}

func (p Digest) Execute(args []string, state storage.State) error {
	config, err := p.ParseArgs(args, state)
	if err != nil {
		return err
	}
	_, err = p.initialize(config, state)

	// choose one manifest from the artifacts galary.

	return err
}

func (p Digest) initialize(config DiegestConfig, state storage.State) (storage.State, error) {
	var err error
	state, err = p.envIDManager.Sync(state, config.Name)
	if err != nil {
		return storage.State{}, fmt.Errorf("Env id manager sync: %s", err)
	}

	err = p.stateStore.Set(state)
	if err != nil {
		return storage.State{}, fmt.Errorf("Save state: %s", err)
	}

	artifactsVarsFilePath, err := p.stateStore.GetMainArtifactVarsFilePath()
	if err != nil {
		return storage.State{}, err
	}

	if err := digesterApis.Run(state, artifactsVarsFilePath); err != nil {
		return storage.State{}, fmt.Errorf("Call digester: %s", err)
	}

	err = p.pickUpManifest()

	return state, err
}

func combine(is digesterCommon.InfraSize, pl digesterCommon.ProgrammingLanguage) string {
	return string(is) + "|" + string(pl)
}

// TODO(zhongyi) This is a stub, should pick the manifest up based on the q&a file.
func (p Digest) pickUpManifest() error {
	artifactsVarsFilePath, err := p.stateStore.GetMainArtifactVarsFilePath()
	if err != nil {
		return err
	}

	artifactFilePath, err := p.stateStore.GetMainArtifactFilePath()
	if err != nil {
		return err
	}

	bp, err := digesterApis.ImportBlueprintYaml(artifactsVarsFilePath)
	if err != nil {
		return err
	}

	var manifestFilePath string
	switch combine(bp.Infra.Size, bp.NonInfra.ProgrammingLanguage) {
	case combine(digesterCommon.SizeSmall, digesterCommon.PHP):
		manifestFilePath = "/manifests/small_php.yml"
	case combine(digesterCommon.SizeMedium, digesterCommon.PHP):
		manifestFilePath = "/manifests/medium_php.yml"
	case combine(digesterCommon.SizeLarge, digesterCommon.PHP):
		manifestFilePath = "/manifests/large_php.yml"
	case combine(digesterCommon.SizeMaximum, digesterCommon.PHP):
		manifestFilePath = "/manifests/maximum_php.yml"
	default:
		return fmt.Errorf("we only support php")

	}

	content, err := builtinmanifests.FSByte(false, manifestFilePath)
	if err != nil {
		return err
	}
	err = p.fs.WriteFile(artifactFilePath, content, 0644)
	return err
}
