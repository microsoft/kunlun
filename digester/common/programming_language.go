package common

import (
	"fmt"
	"strings"
)

type ProgrammingLanguage string

const (
	UnknownProgrammingLanguage ProgrammingLanguage = "unknown"
	PHP                        ProgrammingLanguage = "php"
)

func ParseProgrammingLanguage(pl string) (ProgrammingLanguage, error) {
	pl = strings.ToLower(pl)
	switch pl {
	case "php":
		return PHP, nil
	default:
		return UnknownProgrammingLanguage, fmt.Errorf("Not support the language")
	}
}
