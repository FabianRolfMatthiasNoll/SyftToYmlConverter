package internal

import (
	"encoding/json"
	"log"
	"os"
)

type Syft struct {
	Artifacts []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Version   string `json:"version"`
		Type      string `json:"type"`
		FoundBy   string `json:"foundBy"`
		Locations []struct {
			Path string `json:"path"`
		} `json:"locations"`
		//Licenses     []interface{} `json:"licenses"`
		Licenses     []string `json:"licenses"`
		Language     string   `json:"language"`
		Cpes         []string `json:"cpes"`
		Purl         string   `json:"purl"`
		MetadataType string   `json:"metadataType"`
		Metadata     struct {
		} `json:"metadata"`
	} `json:"artifacts"`
	ArtifactRelationships []struct {
		Parent string `json:"parent"`
		Child  string `json:"child"`
		Type   string `json:"type"`
	} `json:"artifactRelationships"`
	Source struct {
		ID     string `json:"id"`
		Type   string `json:"type"`
		Target string `json:"target"`
	} `json:"source"`
	Schema struct {
		Version string `json:"version"`
		URL     string `json:"url"`
	} `json:"schema"`
}

func (syft *Syft) ReadJson(path string) (*Syft, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(file, syft)
	return syft, err
}
