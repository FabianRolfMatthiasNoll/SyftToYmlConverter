package main

import (
	"log"
	"strings"
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/handler"
)

func main() {
	syftReader := internal.Syft{}
	var manager *Manager
	syft, syftErr := syftReader.ReadJson("../testfiles/dependencies_angular.json")

	if syftErr != nil {
		log.Print(syftErr)
	}

	switch {
	case strings.Contains(syft.Artifacts[0].Purl, "dotnet"):
		manager = NewManager(handler.Dotnet{})
	case strings.Contains(syft.Artifacts[0].Purl, "golang"):
		manager = NewManager(handler.Go{})
	case strings.Contains(syft.Artifacts[0].Purl, "npm"):
		manager = NewManager(handler.Npm{})
	}

	err := manager.Run(syft)
	if err != nil {
		log.Println(err)
	}
}
