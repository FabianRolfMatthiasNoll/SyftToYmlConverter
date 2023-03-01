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
		Licenses     []interface{} `json:"licenses"`
		Language     string        `json:"language"`
		Cpes         []string      `json:"cpes"`
		Purl         string        `json:"purl"`
		MetadataType string        `json:"metadataType"`
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
	Distro struct {
	} `json:"distro"`
	Descriptor struct {
		Name          string `json:"name"`
		Version       string `json:"version"`
		Configuration struct {
			ConfigPath         string   `json:"configPath"`
			Verbosity          int      `json:"verbosity"`
			Quiet              bool     `json:"quiet"`
			Output             []string `json:"output"`
			OutputTemplatePath string   `json:"output-template-path"`
			File               string   `json:"file"`
			CheckForAppUpdate  bool     `json:"check-for-app-update"`
			Dev                struct {
				ProfileCPU bool `json:"profile-cpu"`
				ProfileMem bool `json:"profile-mem"`
			} `json:"dev"`
			Log struct {
				Structured   bool   `json:"structured"`
				Level        string `json:"level"`
				FileLocation string `json:"file-location"`
			} `json:"log"`
			Catalogers interface{} `json:"catalogers"`
			Package    struct {
				Cataloger struct {
					Enabled bool   `json:"enabled"`
					Scope   string `json:"scope"`
				} `json:"cataloger"`
				SearchUnindexedArchives bool `json:"search-unindexed-archives"`
				SearchIndexedArchives   bool `json:"search-indexed-archives"`
			} `json:"package"`
			Attest struct {
			} `json:"attest"`
			FileMetadata struct {
				Cataloger struct {
					Enabled bool   `json:"enabled"`
					Scope   string `json:"scope"`
				} `json:"cataloger"`
				Digests []string `json:"digests"`
			} `json:"file-metadata"`
			FileClassification struct {
				Cataloger struct {
					Enabled bool   `json:"enabled"`
					Scope   string `json:"scope"`
				} `json:"cataloger"`
			} `json:"file-classification"`
			FileContents struct {
				Cataloger struct {
					Enabled bool   `json:"enabled"`
					Scope   string `json:"scope"`
				} `json:"cataloger"`
				SkipFilesAboveSize int           `json:"skip-files-above-size"`
				Globs              []interface{} `json:"globs"`
			} `json:"file-contents"`
			Secrets struct {
				Cataloger struct {
					Enabled bool   `json:"enabled"`
					Scope   string `json:"scope"`
				} `json:"cataloger"`
				AdditionalPatterns struct {
				} `json:"additional-patterns"`
				ExcludePatternNames []interface{} `json:"exclude-pattern-names"`
				RevealValues        bool          `json:"reveal-values"`
				SkipFilesAboveSize  int           `json:"skip-files-above-size"`
			} `json:"secrets"`
			Registry struct {
				InsecureSkipTLSVerify bool          `json:"insecure-skip-tls-verify"`
				InsecureUseHTTP       bool          `json:"insecure-use-http"`
				Auth                  []interface{} `json:"auth"`
			} `json:"registry"`
			Exclude     []interface{} `json:"exclude"`
			Platform    string        `json:"platform"`
			Name        string        `json:"name"`
			Parallelism int           `json:"parallelism"`
		} `json:"configuration"`
	} `json:"descriptor"`
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
