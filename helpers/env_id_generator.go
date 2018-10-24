package helpers

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"time"
)

type EnvIDGenerator struct {
	reader io.Reader
}

func NewEnvIDGenerator(reader io.Reader) EnvIDGenerator {
	return EnvIDGenerator{
		reader: reader,
	}
}

func (e EnvIDGenerator) Generate() (string, error) {
	lake, err := e.randomLake()
	if err != nil {
		return "", err
	}
	timestamp := time.Now().UTC().Format("2006-01-02t15-04z")

	return fmt.Sprintf("kun-lun-env-%s-%s", lake, timestamp), nil
}

func (e EnvIDGenerator) randomLake() (string, error) {
	lakes := []string{
		"shanghai",
		"nanjing",
		"beijing",
	}

	lakeIdx, err := rand.Int(e.reader, big.NewInt(int64(len(lakes))))
	if err != nil {
		return "", err
	}
	return lakes[lakeIdx.Int64()], nil
}
