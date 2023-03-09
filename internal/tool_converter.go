package internal

import (
	"sort"
	"strings"
	"syfttoymlconverter/internal/model"
)

func ToolToDependencies(syft *Syft) *model.SBOM {
	result := &model.SBOM{}

	for _, d := range syft.Artifacts {
		//Generating Language out of the start of the purl
		language := strings.Split(d.Purl, "/")[0][4:]

		if !contains(result.Languages, language) {
			result.Languages = append(result.Languages, language)
		}
		// TopLevel is really only important for docker and angular still have to think of the way
		// var toplevel = false
		// for _, l := range syft.ArtifactRelationships {
		// 	if l.Parent == d.ID {
		// 		toplevel = true
		// 	}
		// }

		result.Dependencies = append(result.Dependencies, model.Dependencies{
			ImportName: d.Name,
			Language:   language,
			Version:    d.Version,
			Licenses:   d.Licenses,
			//Not sure if we should keep purl for tool independency
			Purl: d.Purl,
			ID:   d.ID,
		})
	}
	sortDependenciesByLanguage(result)
	//sortLanguages(result)
	return result
}

func sortDependenciesByLanguage(sbom *model.SBOM) {
	sort.Slice(sbom.Dependencies, func(i, j int) bool {
		return sbom.Dependencies[i].Language < sbom.Dependencies[j].Language
	})
}

// func sortLanguages(sbom *model.SBOM) {
// 	sort.Slice(sbom.Languages, func(i, j int) bool {
// 		return sbom.Languages[i] < sbom.Languages[j]
// 	})
// }

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
