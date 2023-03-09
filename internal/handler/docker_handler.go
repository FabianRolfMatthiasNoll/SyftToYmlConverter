package handler

import (
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/model"
)

type Docker struct{}

func (Docker) FetchMetadata(syft *internal.Syft) (model.BuildInfo, error) {
	return model.BuildInfo{}, nil
}
