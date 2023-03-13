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

// Document is the main structure for generating the files
type Document struct {
	Header              Header              `validate:"required" yaml:"header"`
	FrontPage           FrontPage           `validate:"required" yaml:"frontPage"`
	ReferencedDocuments ReferencedDocuments `validate:"required" yaml:"referenceDocuments"`
	Libraries           []Library           `validate:"required" yaml:"libraries"`
}

// Header has the specified information about the document itself
type Header struct {
	Title          string `validate:"required" yaml:"title"`
	DocumentName   string `validate:"required" yaml:"documentName"`
	DocumentNumber string `validate:"required" yaml:"documentNumber"`
	DocVersion     string `validate:"required" yaml:"docversion"`
}

// FrontPage shows mandatory information at the first page
type FrontPage struct {
	Creator   Person     `validate:"required" yaml:"creator"`
	Reviewers []Reviewer `validate:"required,dive" yaml:"reviewers"`
	Approver  Person     `validate:"required" yaml:"approver"`
	Histories []History  `validate:"required,dive" yaml:"histories"`
}

// Reviewer structure
type Reviewer struct {
	Reviewer Person `validate:"required" yaml:"reviewer"`
}

// Person described a contributor of the document
type Person struct {
	Name       string `validate:"required" yaml:"name"`
	Profession string `validate:"required" yaml:"profession"`
}

// History structure
type History struct {
	History HistoryEntry `validate:"required" yaml:"history"`
}

// HistoryEntry contains informations of previous document versions
type HistoryEntry struct {
	Version                string `validate:"required" yaml:"version"`
	BeginOfValidation      string `validate:"required" yaml:"beginOfValidation"`
	ReasonAndContentColumn string `validate:"required" yaml:"reasonAndContentColumn"`
}

// ReferencedDocuments shows informations about documents which are referenced for another document
type ReferencedDocuments struct {
	Reference      string `validate:"required" yaml:"reference"`
	Description    string `validate:"required" yaml:"description"`
	DocumentNumber string `validate:"required" yaml:"documentNumber"`
}

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
		if len(d.Parents) > 3 {
			continue
		}
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
