package utils

type ErrUnsupportedRfapVersion struct{}

func (e *ErrUnsupportedRfapVersion) Error() string {
	return "Unsupported rfap version"
}

type ErrIsDir struct{}

func (e *ErrIsDir) Error() string {
	return "Is a directory"
}

type ErrIsNotDir struct{}

func (e *ErrIsNotDir) Error() string {
	return "Is not a directory"
}

type ErrSetReadTimeoutFailed struct{}

func (e *ErrSetReadTimeoutFailed) Error() string {
	return "Failed to set read timeout"
}

type ErrClientCrashed struct{}

func (e *ErrClientCrashed) Error() string {
	return "Client crashed"
}

type ErrAccessDenied struct{}

func (e *ErrAccessDenied) Error() string {
	return "Data in different packets doesn't match"
}

type ErrInvalidContentLength struct{}

func (e *ErrInvalidContentLength) Error() string {
	return "Invalid content length"
}

type ErrChecksumsNotMatching struct{}

func (e *ErrChecksumsNotMatching) Error() string {
	return "Checkums don't match"
}

type ErrCalculationFailed struct{}

func (e *ErrCalculationFailed) Error() string {
	return "Calculation failed"
}
