package apis

import (
	"path"

	artifacts "github.com/Microsoft/kunlun/artifacts"
	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/Microsoft/kunlun/common/storage"
)

type Patching struct {
	stateStore storage.Store
	fs         fileio.Fs
}

func NewPatching(
	stateStore storage.Store,
	fs fileio.Fs,
) Patching {
	return Patching{
		stateStore: stateStore,
		fs:         fs,
	}
}

func (p Patching) ProvisionManifest() (*artifacts.Manifest, error) {
	mainArtifactFilePath, err := p.stateStore.GetMainArtifactFilePath()

	if err != nil {
		return nil, err
	}
	content, err := p.fs.ReadFile(mainArtifactFilePath)
	template := NewTemplate(content)

	// construct the ops
	artifactsPatchDir, err := p.stateStore.GetArtifactsPatchDir()
	fileInfos, err := p.fs.ReadDir(artifactsPatchDir)
	opsFileArgs := []OpsFileArg{}
	for _, fileInfo := range fileInfos {
		fileArg := OpsFileArg{
			fileReader: p.fs,
		}
		fileArg.UnmarshalFlag(path.Join(artifactsPatchDir, fileInfo.Name()))
		opsFileArgs = append(opsFileArgs, fileArg)
	}

	opsFlags := OpsFlags{
		OpsFiles: opsFileArgs,
	}

	// build the artifact vars file
	artifactVarsFilePath, err := p.stateStore.GetMainArtifactVarsFilePath()
	varsFileArg := VarsFileArg{
		fs: p.fs,
	}
	err = varsFileArg.UnmarshalFlag(artifactVarsFilePath)
	if err != nil {
		return nil, err
	}

	varsStore := VarsFSStore{
		fs: p.fs,
	}

	varsStoreFilePath, err := p.stateStore.GetMainArtifactVarsStoreFilePath()
	if err != nil {
		return nil, err
	}

	err = varsStore.UnmarshalFlag(varsStoreFilePath)

	if err != nil {
		return nil, err
	}

	varsFlags := VarFlags{
		VarsFiles:   []VarsFileArg{varsFileArg},
		VarsFSStore: varsStore,
	}
	evalOpts := EvaluateOpts{
		ExpectAllKeys:     true,
		ExpectAllVarsUsed: false,
	}
	content, err = template.Evaluate(varsFlags.AsVariables(), opsFlags.AsOp(), evalOpts)
	if err != nil {
		return nil, err
	}
	manifest, err := artifacts.NewManifestFromYAML(content)
	return manifest, err
}
