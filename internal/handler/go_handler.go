package handler

import (
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/api_interfaces"
	"syfttoymlconverter/internal/model"
)

type Go struct{}

func (Go) FetchMetadata(syft *internal.Syft) (model.BuildInfo, error) {

	models, _ := api_interfaces.ParseEmbeddedModules(syft)

	api_interfaces.SetRepoInfoPooled(&models, 5)

	return models, nil
}
