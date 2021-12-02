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

type UnsupportedRfapVersionError struct {
}

func (e *UnsupportedRfapVersionError) Error() string {
	return "Unsupported rfap version"
}
