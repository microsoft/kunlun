package apis

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Microsoft/kunlun/common/fileio"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

type VarsFSStore struct {
	fs                    fileio.Fs
	ValueGeneratorFactory ValueGeneratorFactory

	path string
}

func NewVarsFSStore(fs fileio.Fs) VarsFSStore {
	return VarsFSStore{
		fs: fs,
	}
}

var _ Variables = VarsFSStore{}

func (s VarsFSStore) IsSet() bool { return len(s.path) > 0 }

func (s VarsFSStore) Get(varDef VariableDefinition) (interface{}, bool, error) {
	vars, err := s.load()
	if err != nil {
		return nil, false, err
	}

	val, found := vars[varDef.Name]
	if found {
		return val, true, nil
	}

	if len(varDef.Type) == 0 {
		return nil, false, nil
	}

	val, err = s.generateAndSet(varDef)
	if err != nil {
		return nil, false, errors.New(err.Error() + fmt.Sprintf("Generating variable '%s'", varDef.Name))
	}

	return val, true, nil
}

func (s VarsFSStore) List() ([]VariableDefinition, error) {
	vars, err := s.load()
	if err != nil {
		return nil, err
	}

	return vars.List()
}

func (s VarsFSStore) generateAndSet(varDef VariableDefinition) (interface{}, error) {
	generator, err := s.ValueGeneratorFactory.GetGenerator(varDef.Type)
	if err != nil {
		return nil, err
	}

	val, err := generator.Generate(varDef.Options)
	if err != nil {
		return nil, err
	}

	err = s.set(varDef.Name, val)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (s VarsFSStore) set(key string, val interface{}) error {
	vars, err := s.load()
	if err != nil {
		return err
	}

	vars[key] = val

	return s.save(vars)
}

func (s VarsFSStore) load() (StaticVariables, error) {
	vars := StaticVariables{}

	// file exists?
	_, err := s.fs.Stat(s.path)

	if err != nil {
		if !strings.Contains(err.Error(), "no such file") {
			return StaticVariables{}, err
		}
	}
	if err == nil {
		bytes, err := s.fs.ReadFile(s.path)
		if err != nil {
			return vars, err
		}

		err = yaml.Unmarshal(bytes, &vars)
		if err != nil {
			return vars, errors.New(err.Error() + fmt.Sprintf("Deserializing variables file store '%s'", s.path))
		}
	}
	if vars == nil {
		return StaticVariables{}, nil
	}

	return vars, nil
}

func (s VarsFSStore) save(vars StaticVariables) error {
	bytes, err := yaml.Marshal(vars)
	if err != nil {
		return errors.New(err.Error() + ("Serializing variables"))
	}

	err = s.fs.WriteFile(s.path, bytes, 0644)
	if err != nil {
		return errors.New(err.Error() + fmt.Sprintf("Writing variables to file store '%s'", s.path))
	}

	return nil
}

func (s *VarsFSStore) UnmarshalFlag(data string) error {
	if len(data) == 0 {
		return errors.New("Expected file path to be non-empty")
	}

	// TODO add support for the ~ home directory.
	absPath, err := filepath.Abs(data)
	if err != nil {
		return errors.New(err.Error() + fmt.Sprintf("Getting absolute path '%s'", data))
	}

	(*s).path = absPath
	(*s).ValueGeneratorFactory = NewValueGeneratorConcrete(nil)

	return nil
}
