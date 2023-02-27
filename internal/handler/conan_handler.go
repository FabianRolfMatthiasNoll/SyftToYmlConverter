package handler

import (
	"fmt"
	"log"
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/api_interfaces"
	"syfttoymlconverter/internal/model"
)

type Conan struct{}

func (Conan) FetchMetadata(syft *internal.Syft) (model.BuildInfo, error) {
	conanInfo, err := api_interfaces.GetMetadata("zlib", "v2.3.1")
	if err != nil {
		log.Printf("Error getting metadata: %v", err)
	}
	conan := api_interfaces.ParseConanOutput(string(conanInfo))
	fmt.Println(conan)
	return model.BuildInfo{}, nil
}
