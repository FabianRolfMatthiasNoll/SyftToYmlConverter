package api_interfaces

import (
	"fmt"
	"strings"
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/model"
	"syfttoymlconverter/internal/provider"

	"github.com/gammazero/workerpool"
	"github.com/rs/zerolog/log"
)

const (
	answer1fmt = "Manufacturer: %s\nVersion: %s\nLicense: %s\n\n%s"
	answer2    = "The hardware specification is documented in [DPS].\nThe software specification is outlined in [SRS]."
	answer3    = "Software Requirements are captured in [SRS]. Software Tests are described in [SVP] and [VTP]. Traceability is ensured by [EVDR].\n\n(1) The OTS SW was incorporated in the device during installation of the system. There is no possibility to see, remove or change the system files."
	answer4Fmt = "The software is needed as dependency of %s.\n\nThere are no specialized requirements defined for this component. Requirements for the system are specified in [SRS].\n\nThe OTS SW does not link with software outside the system."
	answer5    = "Software Tests are described in [SVP] and [VTP].\n\nThe OTS SW is incorporated in the device during installation of the system. There is no possibility to see, remove or change the system files."
	answer6    = "The OTS SW is incorporated in the device during installation of the system and it will be ensured, that the user can not see, remove or change the system files.\nConfiguration and Version of the OTS is kept under version control in Git\n\nThe lifecycle of the OTS will be maintained using the [FOSS] process."
)

type Go struct {
	Path    string // Import path, such as "github.com/mitchellh/golicense"
	SubPath string // matches the trailing import version specifiers like `/v12`
	Version string // Version like "v1.2.3"
	Hash    string // Hash such as "h1:abcd1234"
}

func SyftToModule(syft *internal.Syft) ([]Go, error) {
	var result []Go
	for i, data := range syft.Artifacts {
		//scanning the binary of a go file returns as first element itself
		if i == 0 {
			continue
		}
		data.Name = removeExtraPath(data.Name)
		next := Go{
			Path:    data.Name,
			SubPath: data.Version[:strings.Index(data.Version, ".")],
			Version: data.Version,
			Hash:    data.ID,
		}

		result = append(result, next)
	}
	return result, nil
}

func ParseEmbeddedModules(syft *internal.Syft) (model.BuildInfo, error) {

	modules, err := SyftToModule(syft)
	if err != nil {
		return model.BuildInfo{}, err
	}

	return model.BuildInfo{
		Path:    syft.Artifacts[0].Locations[0].Path,
		Mod:     "Mod",
		Modules: toModule(modules),
	}, nil
}

func toModule(modules []Go) []model.Module {
	var result []model.Module

	for _, m := range modules {
		result = append(result, model.Module{
			Path:    m.Path,
			SubPath: m.SubPath,
			Version: m.Version,
			Hash:    m.Hash,
		})
	}

	return result
}

func SetRepoInfoPooled(info *model.BuildInfo, workers int) {
	wp := workerpool.New(workers)

	for i := range info.Modules {
		module := &info.Modules[i]

		wp.Submit(func() {
			log.Info().Msgf("fetching module info for %s", module.String())

			module.Info = provider.FetchModuleInfo(module.Path, module.Version)
		})
	}

	wp.StopWait()
}

// This will be deprecated in the future because we will use the Part of the FOSSer Tool
func ModelToLibrary(info *model.BuildInfo) model.Librarys {
	libs := model.Librarys{Libraries: []model.Library{}}
	for _, d := range info.Modules {
		var lib model.Library
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

func removeExtraPath(input string) string {
	parts := strings.Split(input, "/")
	if len(parts) >= 3 {
		parts = parts[:3]
	}
	return strings.Join(parts, "/")
}
