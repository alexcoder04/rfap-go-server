package main

type HeaderMetadata struct {
	ErrorCode    int    `yaml:"ErrorCode"`
	ErrorMessage string `yaml:"ErrorMessage"`
	Path         string `yaml:"Path"`
	Type         string `yaml:"Type"`
	Modified     int    `yaml:"Modified"`

	FileSize int    `yaml:"FileSize"`
	FileType string `yaml:"FileType"`

	DirectorySize  int `yaml:"DirectorySize"`
	ElementsNumber int `yaml:"ElementsNumber"`
}
