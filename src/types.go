package main

type HeaderValues struct {
	ErrorCode    int    `yaml:"ErrorCode"`
	ErrorMessage string `Yaml:"ErrorMessage"`
	Path         string `yaml:"FilePath"`
	Type         rune   `yaml:"Type"`
	Modified     int    `yaml:"Modified"`

	FileSize int    `yaml:"FileSize"`
	FileType string `yaml:"FileType"`

	DirectorySize  int `yaml:"DirectorySize"`
	ElementsNumber int `yaml:"ElementsNumber"`
}
