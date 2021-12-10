package main

type HeaderMetadata struct {
	ErrorCode    int    `yaml:"ErrorCode"`
	ErrorMessage string `yaml:"ErrorMessage"`

	RequestDetails []string `yaml:"RequestDetails"`

	Path         string `yaml:"Path"`
	Type         string `yaml:"Type"`
	Modified     int    `yaml:"Modified"`
	PacketsTotal int    `yaml:"PacketsTotal"`
	PacketNumber int    `yaml:"PacketNumber"`

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

type ErrSetReadTimeoutFailed struct{}

func (e *ErrSetReadTimeoutFailed) Error() string {
	return "Failed to set read timeout"
}

type ErrClientCrashed struct{}

func (e *ErrClientCrashed) Error() string {
	return "Client crashed"
}

type ErrInvalidPacketNumber struct{}

func (e *ErrInvalidPacketNumber) Error() string {
	return "Invalid total number of packets"
}

type ErrDifferentPacketsDontMatch struct{}

func (e *ErrDifferentPacketsDontMatch) Error() string {
	return "Data in different packets doesn't match"
}

type ErrAccessDenied struct{}

func (e *ErrAccessDenied) Error() string {
	return "Data in different packets doesn't match"
}
