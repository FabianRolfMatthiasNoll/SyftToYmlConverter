package provider

import (
	"fmt"
	"regexp"
	"strings"

	"syfttoymlconverter/internal/model"

	"github.com/rs/zerolog/log"
)

// is the regexp matching the package for a golang import
var golangRegEx = regexp.MustCompile(`^golang\.org\/x\/(.+)$`)

// TODO: check if pkg.go.dev does provide an API to fetch needed information
// it does a better job regarding license info (see: https://github.com/golang/go/issues/36785)
func FetchModuleInfo(source, version string) model.RepoInfo {
	source = resolve(source)
	tag := strings.TrimSuffix(version, "+incompatible")

	return githubClient.getInfo(source, tag)
}

func FetchLicenseText(source, spdx string) (string, bool) {
	log.Info().Msgf("Fetching license text for %s", source)

	source = resolve(source)

	text, err := githubClient.getLicenseFromRepo(source)
	if err != nil {
		log.Debug().Err(err).Msgf("Failed to get license text for %s", source)

		// no license file in repo, return general license text
		text, err = githubClient.getSpdxLicense(spdx)
		if err != nil {
			log.Debug().Err(err).Msgf("Failed to get spdx data for %s", spdx)

			return "", false
		}
	}

	return text, true
}

func ResolveWebsiteLink(modulePath string) string {
	return fmt.Sprintf("https://pkg.go.dev/%s", modulePath)
}

func ResolveLicenseLink(modulePath string) string {
	return fmt.Sprintf("https://pkg.go.dev/%s?tab=licenses", modulePath)
}

func resolve(path string) string {
	// change source from golang.org to github.com/golang
	result, ok := golang(path)
	if ok {
		return result
	}

	return path
}

func golang(path string) (string, bool) {
	ms := golangRegEx.FindStringSubmatch(path)
	if ms == nil {
		return "", false
	}

	// Matches, convert to github
	p := fmt.Sprintf("github.com/golang/%s", ms[1])

	return p, true
}
