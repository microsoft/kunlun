package errors

type NotSupportedError struct {
}

func (e *NotSupportedError) Error() string {
	return "not supported."
}
