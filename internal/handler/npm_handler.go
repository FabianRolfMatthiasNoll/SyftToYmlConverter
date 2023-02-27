package handler

import (
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/api_interfaces"
	"syfttoymlconverter/internal/model"
)

type Npm struct{}

func (Npm) FetchMetadata(syft *internal.Syft) (model.BuildInfo, error) {
	var npm api_interfaces.NPM
	models, _ := npm.ParseEmbeddedModules(syft)
	npm.SetRepoInfo(syft, &models)
	return models, nil
}
