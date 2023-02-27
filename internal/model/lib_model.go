package model

type Librarys struct {
	Libraries []Library `json:"libraries"`
}

// Library structure
type Library struct {
	Source      string            `validate:"required" yaml:"source"`
	Submodule   string            `yaml:"submodule"`
	Release     string            `validate:"required" yaml:"release"`
	LibraryData TableMainTemplate `validate:"required" yaml:"libraryTable"`
}

// TableMainTemplate contains the content of a library table
type TableMainTemplate struct {
	Manufacturer   string `validate:"required" yaml:"manufacturer"`
	Software       string `validate:"required" yaml:"software"`
	Summary        string `validate:"required" yaml:"summary"`
	Version        string `validate:"required" yaml:"version"`
	License        string `validate:"required" yaml:"license"`
	Function       string `validate:"required" yaml:"function"`
	Incorporated   string `validate:"required" yaml:"incorporated"`
	LevelOfConcern string `validate:"required" yaml:"levelOfConcern"`
	Answer1        string `validate:"required" yaml:"answer1"`
	Answer2        string `validate:"required" yaml:"answer2"`
	Answer3        string `validate:"required" yaml:"answer3"`
	Answer4        string `validate:"required" yaml:"answer4"`
	Answer5        string `validate:"required" yaml:"answer5"`
	Answer6        string `validate:"required" yaml:"answer6"`
}
