package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Microsoft/kunlun/artifacts/builtinmanifests"
	qgraph "github.com/Microsoft/kunlun/artifacts/qgraph"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/flags"
	"github.com/Microsoft/kunlun/common/helpers"
	"github.com/Microsoft/kunlun/common/storage"
	"github.com/Microsoft/kunlun/common/ui"
	"github.com/andyliuliming/gquiz"
	yaml "gopkg.in/yaml.v2"
)

type Digest struct {
	stateStore   storage.Store
	envIDManager helpers.EnvIDManager

	fs fileio.Fs
	ui *ui.UI
}

type DiegestConfig struct {
	Name string
}

func NewDigest(
	stateStore storage.Store,
	envIDManager helpers.EnvIDManager,
	fs fileio.Fs,
	ui *ui.UI,
) Digest {
	return Digest{
		stateStore:   stateStore,
		envIDManager: envIDManager,
		fs:           fs,
		ui:           ui,
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

	var defaults gquiz.QResult
	if _, err := os.Stat(artifactsVarsFilePath); !os.IsNotExist(err) {
		varsContent, err := ioutil.ReadFile(artifactsVarsFilePath)
		if err != nil {
			return storage.State{}, err
		}
		yaml.Unmarshal(varsContent, &defaults)
	}
	qResult, err := p.doQuiz(defaults)
	bpBytes, _ := yaml.Marshal(qResult)
	err = ioutil.WriteFile(artifactsVarsFilePath, bpBytes, 0644)

	content, err := builtinmanifests.FSByte(false, path.Join("/manifests", (qResult)["final_artifact"]))
	if err != nil {
		return state, err
	}

	artifactFilePath, err := p.stateStore.GetMainArtifactFilePath()
	if err != nil {
		return state, err
	}
	err = p.fs.WriteFile(artifactFilePath, content, 0644)
	return state, err
}

func (p Digest) doQuiz(defaults gquiz.QResult) (gquiz.QResult, error) {
	fs := qgraph.FS(false)
	qgraphFolder := "/manifests"
	file, err := fs.Open(qgraphFolder)
	if err != nil {
		return nil, err
	}
	files, err := file.Readdir(0)
	var sb strings.Builder
	for _, f := range files {
		filePath := path.Join(qgraphFolder, f.Name())
		content, err := qgraph.FSByte(false, filePath)
		if err != nil {
			return nil, err
		}
		_, err = sb.Write(content)
		sb.WriteString("\n")
		if err != nil {
			return nil, err
		}
	}
	quizeBuilder := gquiz.QuizBuilder{}

	qGraph, err := quizeBuilder.BuildQGraph([]byte(sb.String()))
	if err != nil {
		return nil, err
	}

	// read the main vars as the default qresult.
	quizExecutor := gquiz.NewQuizExecutor(p.ui, defaults)
	qResult, err := quizExecutor.Execute(&qGraph)
	if err != nil {
		return nil, err
	}
	return qResult, nil
}
