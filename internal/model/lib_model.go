package model

import (
	"fmt"
	"strings"
)

const (
	answer1fmt = "Manufacturer: %s\nVersion: %s\nLicense: %s\n\n%s"
	answer2    = "The hardware specification is documented in [DPS].\nThe software specification is outlined in [SRS]."
	answer3    = "Software Requirements are captured in [SRS]. Software Tests are described in [SVP] and [VTP]. Traceability is ensured by [EVDR].\n\n(1) The OTS SW was incorporated in the device during installation of the system. There is no possibility to see, remove or change the system files."
	answer4Fmt = "The software is needed as dependency of %s.\n\nThere are no specialized requirements defined for this component. Requirements for the system are specified in [SRS].\n\nThe OTS SW does not link with software outside the system."
	answer5    = "Software Tests are described in [SVP] and [VTP].\n\nThe OTS SW is incorporated in the device during installation of the system. There is no possibility to see, remove or change the system files."
	answer6    = "The OTS SW is incorporated in the device during installation of the system and it will be ensured, that the user can not see, remove or change the system files.\nConfiguration and Version of the OTS is kept under version control in Git\n\nThe lifecycle of the OTS will be maintained using the [FOSS] process."
)

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

func ModelToLibrary(info *BuildInfo) Librarys {
	libs := Librarys{Libraries: []Library{}}
	for _, d := range info.Modules {
		var lib Library
		lib.Source = d.Path
		lib.Submodule = d.SubPath

		if !d.Info.Release.IsZero() {
			lib.Release = d.Info.Release.Format("2006-01-02")
		}

		lib.LibraryData.Function = "Library"
		lib.LibraryData.Version = d.Version
		lib.LibraryData.Manufacturer = d.Info.FullName
		lib.LibraryData.Summary = d.Info.Description
		lib.LibraryData.License = d.Info.SPDX

		// use last path as software name: github.com/integrii/flaggy -> flaggy
		s := strings.Split(d.Path, "/")
		lib.LibraryData.Software = s[len(s)-1]

		// default values
		lib.LibraryData.Answer1 = fmt.Sprintf(
			answer1fmt,
			lib.LibraryData.Manufacturer,
			lib.LibraryData.Version,
			lib.LibraryData.License,
			lib.LibraryData.Summary,
		)

		if len(d.Parents) > 0 {
			lib.LibraryData.Answer4 = fmt.Sprintf(answer4Fmt, strings.Join(d.Parents, ", "))
		}

		lib.LibraryData.Answer2 = answer2
		lib.LibraryData.Answer3 = answer3
		lib.LibraryData.Answer5 = answer5
		lib.LibraryData.Answer6 = answer6

		lib.LibraryData.Incorporated = "Yes"
		lib.LibraryData.LevelOfConcern = "Minor"
		lib.LibraryData.Function = "Library"

		libs.Libraries = append(libs.Libraries, lib)
	}
	return libs
}
