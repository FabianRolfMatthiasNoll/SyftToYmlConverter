package main

import (
	"log"
	"os"
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/model"

	"github.com/goccy/go-yaml"
)

type Lang_Interface interface {
	FetchMetadata(syft *internal.Syft) (model.BuildInfo, error)
}

type Manager struct {
	Lang Lang_Interface
}

func NewManager(l Lang_Interface) *Manager {
	return &Manager{
		Lang: l,
	}
}

func (m *Manager) Run(syft *internal.Syft) error {
	var models model.BuildInfo
	libraries := model.Librarys{Libraries: []model.Library{}}

	models, _ = m.Lang.FetchMetadata(syft)

	libraries = model.ModelToLibrary(&models)

	yamlData, yamlErr := yaml.Marshal(&libraries)
	if yamlErr != nil {
		log.Print(yamlErr)
	}

	fileName := "../test.yaml"
	yamlErr = os.WriteFile(fileName, yamlData, 0644)

	if yamlErr != nil {
		log.Print(yamlErr)
	}
	return nil
}
