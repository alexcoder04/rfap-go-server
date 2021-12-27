package utils

type HeaderMetadata struct {
	ErrorCode    int    `yaml:"ErrorCode"`
	ErrorMessage string `yaml:"ErrorMessage"`

	RequestDetails []string `yaml:"RequestDetails"`

	Path        string `yaml:"Path"`
	Type        string `yaml:"Type"`
	Modified    int    `yaml:"Modified"`
	Destination string `yaml:"Destination"`

	FileSize int    `yaml:"FileSize"`
	FileType string `yaml:"FileType"`

	DirectorySize  int `yaml:"DirectorySize"`
	ElementsNumber int `yaml:"ElementsNumber"`
}

type CommandExec func(string) (HeaderMetadata, []byte, error)

type CopySommandExec func(string, string, bool) (HeaderMetadata, []byte, error)

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
