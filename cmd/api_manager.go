package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/api_interfaces"
	"syfttoymlconverter/internal/model"

	"github.com/goccy/go-yaml"
)

type API_Interface interface {
	GetData(string) ([]byte, error)
	ConverseDataToStruct([]byte, string) (model.Library, error)
	CreateAPILink(string) string
}

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

	libraries = api_interfaces.ModelToLibrary(&models)

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

func GetPackageName(purl string) string {
	if strings.Contains(purl, "dotnet") {
		//purl = strings.Replace(purl, "pkg:dotnet", "https://www.nuget.org/packages", -1)
		pos := strings.Index(purl, "@")
		pos2 := strings.LastIndex(purl, "/")
		purl = purl[pos2+1 : pos]
	}
	return purl
}

func DisplayLibraries(libraries []model.Library) {
	for _, l := range libraries {
		fmt.Println("============================================================")
		fmt.Println("Source: ", l.Source)
		fmt.Println("Submodule: ", l.Submodule)
		fmt.Println("Release: ", l.Release)
		fmt.Println("Manufacturer: ", l.LibraryData.Manufacturer)
		fmt.Println("Summary: ", l.LibraryData.Summary)
		fmt.Println("Version: ", l.LibraryData.Version)
		fmt.Println("License: ", l.LibraryData.License)
		fmt.Println("Function: ", l.LibraryData.Function)
		fmt.Println("Incorporated: ", l.LibraryData.Incorporated)
		fmt.Println("LevelOfConcern: ", l.LibraryData.LevelOfConcern)

	}

	// for _, d := range models.Modules {
	// 	fmt.Println("===========================================")
	// 	fmt.Println(d.Version)
	// 	fmt.Println(d.Path)
	// 	fmt.Println(d.SubPath)
	// 	fmt.Println(d.Info.SPDX)
	// 	fmt.Println(d.Info.Description)
	// 	fmt.Println(d.Info.Release)
	// }
}
