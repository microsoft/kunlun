package apis

import "fmt"

type resourceValidator func(m Manifest) error

type validator struct {
	resourceValidators []resourceValidator
}

func (v validator) Validate(m Manifest) error {
	for _, rv := range v.resourceValidators {
		if err := rv(m); err != nil {
			return err
		}
	}
	return nil
}

func newValidator() validator {
	return validator{
		resourceValidators: []resourceValidator{
			vmGroupValidator,
			loadBalancerValidator,
			virtualNetworkValidator,
			mysqlDatabaseValidator,
		},
	}
}

func validationError(errorMessage string, a ...interface{}) error {
	if len(a) > 0 {
		return fmt.Errorf("validation error: %s", fmt.Sprintf(errorMessage, a))
	}
	return fmt.Errorf("validation error: %s", errorMessage)
}
