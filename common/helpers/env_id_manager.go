package helpers

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"time"

	common_errors "github.com/Microsoft/kunlun/common/errors"
	"github.com/Microsoft/kunlun/common/storage"
)

type EnvIDManager struct {
	reader io.Reader
}

func NewEnvIDManager(reader io.Reader) EnvIDManager {
	return EnvIDManager{
		reader: reader,
	}
}

func (e EnvIDManager) Sync(state storage.State, envID string) (storage.State, error) {
	if state.EnvID != "" {
		return state, nil
	}

	err := e.checkFastFail(state.IAAS, envID)
	if err != nil {
		return storage.State{}, err
	}

	if envID == "" {
		state.EnvID, err = e.generate()
		if err != nil {
			return storage.State{}, err
		}
	} else {
		err = e.validateName(envID)
		if err != nil {
			return storage.State{}, err
		}

		state.EnvID = envID
	}

	return state, nil
}

func (e EnvIDManager) generate() (string, error) {
	lake, err := e.randomLake()
	if err != nil {
		return "", err
	}
	timestamp := time.Now().UTC().Format("2006-01-02t15-04z")

	return fmt.Sprintf("kl-env-%s-%s", lake, timestamp), nil
}

func (e EnvIDManager) randomLake() (string, error) {
	lakes := []string{
		"beijing",
		"shanghai",
		"nanjing",
		"qinghai",
	}

	lakeIdx, err := rand.Int(e.reader, big.NewInt(int64(len(lakes))))
	if err != nil {
		return "", err
	}
	return lakes[lakeIdx.Int64()], nil
}

func (e EnvIDManager) checkFastFail(iaas, envID string) error {
	switch iaas {
	case "azure":
		return nil
	case "azurestack":
		return &common_errors.NotImplementedError{}
	default:
		return fmt.Errorf("iaas: %s not supported", iaas)
	}
}

func (e EnvIDManager) validateName(envID string) error {
	matched, err := regexp.MatchString("^(?:[a-z](?:[-a-z0-9]*[a-z0-9])?)$", envID)
	if err != nil {
		return err
	}

	if !matched {
		return errors.New("Names must start with a letter, be all lowercase, and be alphanumeric or hyphenated")
	}

	return nil
}
