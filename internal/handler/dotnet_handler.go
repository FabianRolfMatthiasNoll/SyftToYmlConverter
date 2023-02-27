package handler

import (
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/api_interfaces"
	"syfttoymlconverter/internal/model"
)

type Dotnet struct{}

func (Dotnet) FetchMetadata(syft *internal.Syft) (model.BuildInfo, error) {
	var nuget api_interfaces.Nuget
	models, _ := nuget.ParseEmbeddedModules(syft)
	nuget.SetRepoInfo(syft, &models)
	return models, nil
}
