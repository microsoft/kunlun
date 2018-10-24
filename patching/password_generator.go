package apis

import (
	"crypto/rand"
	"errors"
	"math/big"
)

type passwordGenerator struct{}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

const DefaultPasswordLength = 20

type passwordParams struct {
	Length int `yaml:"length"`
}

var supportedPasswordParams = []string{
	"length",
}

func NewPasswordGenerator() ValueGenerator {
	return passwordGenerator{}
}

func (passwordGenerator) Generate(parameters interface{}) (interface{}, error) {
	var params passwordParams
	err := objToStruct(parameters, &params, supportedPasswordParams)
	if err != nil {
		return nil, errors.New(err.Error() + "Failed to generate password, parameters are invalid")
	}

	if params.Length < 0 {
		return nil, errors.New("Failed to generate password, 'length' param cannot be negative")
	}

	if params.Length == 0 {
		params.Length = DefaultPasswordLength
	}

	lengthLetterRunes := big.NewInt(int64(len(letterRunes)))
	passwordRunes := make([]rune, params.Length)

	for i := range passwordRunes {
		index, err := rand.Int(rand.Reader, lengthLetterRunes)
		if err != nil {
			return nil, err
		}

		passwordRunes[i] = letterRunes[index.Int64()]
	}

	return string(passwordRunes), nil
}
