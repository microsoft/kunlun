package config

import (
	"fmt"
	"path/filepath"

	"github.com/Microsoft/kunlun/common/storage"
)

type Merger struct {
	fs fs
}

func NewMerger(fs fs) Merger {
	return Merger{fs: fs}
}

func (m Merger) MergeGlobalFlagsToState(globalFlags GlobalFlags, state storage.State) (storage.State, error) {
	if globalFlags.IAAS != "" {
		if state.IAAS != "" && globalFlags.IAAS != state.IAAS {
			return storage.State{}, fmt.Errorf("The iaas type cannot be changed for an existing environment. The current iaas type is %s.", state.IAAS)
		}
		state.IAAS = globalFlags.IAAS
	}

	switch state.IAAS {
	case "azure":
		return m.updateAzureState(globalFlags, state)
	}

	return state, nil
}

func copyFlagToState(source string, sink *string) {
	if source != "" {
		*sink = source
	}
}

func copyFlagToStateWithDefault(source string, sink *string, def string) {
	if source == "" {
		*sink = def
	} else {
		*sink = source
	}
}

func (m Merger) updateAzureState(globalFlags GlobalFlags, state storage.State) (storage.State, error) {
	copyFlagToState(globalFlags.AzureRegion, &state.Azure.Region)
	copyFlagToState(globalFlags.AzureEnvironment, &state.Azure.Environment)
	copyFlagToState(globalFlags.AzureSubscriptionID, &state.Azure.SubscriptionID)
	copyFlagToState(globalFlags.AzureTenantID, &state.Azure.TenantID)
	copyFlagToState(globalFlags.AzureClientID, &state.Azure.ClientID)
	copyFlagToState(globalFlags.AzureClientSecret, &state.Azure.ClientSecret)

	return state, nil
}

func (m Merger) readKey(path string) (string, string, error) {
	keyBytes, err := m.fs.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("Reading key: %v", err)
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", "", fmt.Errorf("Getting absolute path to key: %v", err)
	}
	return absPath, string(keyBytes), nil
}
