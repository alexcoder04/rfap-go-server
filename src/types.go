package main

type HeaderMetadata struct {
	ErrorCode    int    `yaml:"ErrorCode"`
	ErrorMessage string `yaml:"ErrorMessage"`

	RequestDetails []string `yaml:"RequestDetails"`

	Path     string `yaml:"Path"`
	Type     string `yaml:"Type"`
	Modified int    `yaml:"Modified"`

	FileSize int    `yaml:"FileSize"`
	FileType string `yaml:"FileType"`

	DirectorySize  int `yaml:"DirectorySize"`
	ElementsNumber int `yaml:"ElementsNumber"`
}

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
