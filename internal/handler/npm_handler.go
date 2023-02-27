package handler

import (
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/model"
)

type Npm struct{}

func (Npm) FetchMetadata(syft *internal.Syft) (model.BuildInfo, error) {
	return model.BuildInfo{}, nil
}
