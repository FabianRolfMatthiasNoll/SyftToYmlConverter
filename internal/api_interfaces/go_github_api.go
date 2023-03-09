package api_interfaces

import (
	"strings"
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/model"
	"syfttoymlconverter/internal/provider"

	"github.com/gammazero/workerpool"
	"github.com/rs/zerolog/log"
)

type Go struct {
	Path    string // Import path, such as "github.com/mitchellh/golicense"
	SubPath string // matches the trailing import version specifiers like `/v12`
	Version string // Version like "v1.2.3"
	Hash    string // Hash such as "h1:abcd1234"
}

func SyftToModule(syft *internal.Syft) ([]Go, error) {
	var result []Go
	for _, data := range syft.Artifacts {
		// //scanning the binary of a go file returns as first element itself
		// if i == 0 {
		// 	continue
		// }
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

func removeExtraPath(input string) string {
	parts := strings.Split(input, "/")
	if len(parts) >= 3 {
		parts = parts[:3]
	}
	return strings.Join(parts, "/")
}
