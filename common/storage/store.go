package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"

	"github.com/Microsoft/kunlun/common/fileio"
	uuid "github.com/nu7hatch/gouuid"
)

const (
	STATE_SCHEMA = 1
	STATE_FILE   = "kl-state.json"
)

type Store struct {
	dir         string
	fs          fs
	stateSchema int
}

type fs interface {
	fileio.FileWriter
	fileio.Remover
	fileio.AllRemover
	fileio.Stater
	fileio.AllMkdirer
	fileio.DirReader
}

func NewStore(dir string, fs fs) Store {
	return Store{
		dir:         dir,
		fs:          fs,
		stateSchema: STATE_SCHEMA,
	}
}

func (s Store) Set(state State) error {
	_, err := s.fs.Stat(s.dir)
	if err != nil {
		return fmt.Errorf("Stat state dir: %s", err)
	}

	if reflect.DeepEqual(state, State{}) {
		return nil
	}

	state.Version = s.stateSchema

	if state.ID == "" {
		uuid, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("Create state ID: %s", err)
		}
		state.ID = uuid.String()
	}

	jsonData, err := json.MarshalIndent(state, "", "\t")
	if err != nil {
		return err
	}

	stateFile := filepath.Join(s.dir, STATE_FILE)
	err = s.fs.WriteFile(stateFile, jsonData, os.FileMode(0644))
	if err != nil {
		return err
	}

	return nil
}

func (s Store) GetStateDir() string {
	return s.dir
}

func (s Store) GetVarsDir() (string, error) {
	return s.getDir("vars", StateMode)
}

// GetArtifactsDir get artifacts folder
func (s Store) GetArtifactsDir() (string, error) {
	return s.getDir("artifacts", os.ModePerm)
}

// GetMainArtifactFilePath get artifacts main file path
func (s Store) GetMainArtifactFilePath() (string, error) {
	artifactsDir, err := s.GetArtifactsDir()
	if err != nil {
		return "", err
	}
	return path.Join(artifactsDir, "main.yml"), nil
}

func (s Store) GetAdminSSHPrivateKeyPath() (string, error) {
	varsFolder, err := s.GetVarsDir()
	if err != nil {
		return "", err
	}
	sshPrivateKeyPath := path.Join(varsFolder, "admin_ssh_private_key")
	return sshPrivateKeyPath, nil
}

// GetMainArtifactVarsFilePath get the variables file path
func (s Store) GetMainArtifactVarsFilePath() (string, error) {
	artifactsDir, err := s.GetVarsDir()
	if err != nil {
		return "", err
	}
	return path.Join(artifactsDir, "main-vars-file.yml"), nil
}

// GetMainArtifactVarsStoreFilePath get the vars store, to store the vars generated.
func (s Store) GetMainArtifactVarsStoreFilePath() (string, error) {
	artifactsDir, err := s.GetVarsDir()
	if err != nil {
		return "", err
	}
	return path.Join(artifactsDir, "main-vars-store.yml"), nil
}

// GetArtifactsPatchDir get the patches folder
func (s Store) GetArtifactsPatchDir() (string, error) {
	return s.getDir("artifacts/patches", os.ModePerm)
}

// GetInfraDir get the infrastructure folder.
func (s Store) GetInfraDir() (string, error) {
	return s.getDir("infra", os.ModePerm)
}

// GetTerraformDir get the terraform folder, this should be the sub folder of infra.
func (s Store) GetTerraformDir() (string, error) {
	return s.getDir("infra/terraform", os.ModePerm)
}

func (s Store) GetDeploymentsDir() (string, error) {
	return s.getDir("deployments", os.ModePerm)
}

func (s Store) GetDeploymentScriptFile() (string, error) {
	deploymentsDir, err := s.GetDeploymentsDir()
	if err != nil {
		return "", err
	}
	return path.Join(deploymentsDir, "deploy.sh"), nil
}

func (s Store) GetAnsibleDir() (string, error) {
	return s.getDir("deployments/ansible", os.ModePerm)
}

func (s Store) GetAnsibleMainFile() (string, error) {
	ansibleDir, err := s.GetAnsibleDir()
	if err != nil {
		return "", err
	}
	return path.Join(ansibleDir, "main.yml"), nil
}

// TODO think about merge the vars dir with the global vars dir.
func (s Store) GetAnsibleVarsDir() (string, error) {
	return s.getDir("deployments/ansible/vars", os.ModePerm)
}

func (s Store) GetAnsibleInventoriesDir() (string, error) {
	return s.getDir("deployments/ansible/inventories", os.ModePerm)
}

func (s Store) getDir(name string, perm os.FileMode) (string, error) {
	dir := filepath.Join(s.dir, name)
	err := s.fs.MkdirAll(dir, perm)
	if err != nil {
		return "", fmt.Errorf("Get %s dir: %s", name, err)
	}
	return dir, nil
}
