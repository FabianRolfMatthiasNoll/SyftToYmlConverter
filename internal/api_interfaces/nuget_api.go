package api_interfaces

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"syfttoymlconverter/internal"
	"syfttoymlconverter/internal/model"
	"time"

	"github.com/TwiN/go-color"
)

type Module struct {
	Path    string // Import path, such as "github.com/mitchellh/golicense"
	SubPath string // matches the trailing import version specifiers like `/v12`
	Version string // Version like "v1.2.3"
	Hash    string // Hash such as "h1:abcd1234"
}

// Structure of NUGET API Call https://api.nuget.org/v3/registration5-semver1/{PackageNameLowerCase}/index.json
type Nuget struct {
	ID              string    `json:"@id"`
	Type            []string  `json:"@type"`
	CommitID        string    `json:"commitId"`
	CommitTimeStamp time.Time `json:"commitTimeStamp"`
	Count           int       `json:"count"`
	Items           []struct {
		ID              string    `json:"@id"`
		Type            string    `json:"@type"`
		CommitID        string    `json:"commitId"`
		CommitTimeStamp time.Time `json:"commitTimeStamp"`
		Count           int       `json:"count"`
		Items           []struct {
			ID              string    `json:"@id"`
			Type            string    `json:"@type"`
			CommitID        string    `json:"commitId"`
			CommitTimeStamp time.Time `json:"commitTimeStamp"`
			CatalogEntry    struct {
				IDAT                     string    `json:"@id"`
				Type                     string    `json:"@type"`
				Authors                  string    `json:"authors"`
				Description              string    `json:"description"`
				IconURL                  string    `json:"iconUrl"`
				ID                       string    `json:"id"`
				Language                 string    `json:"language"`
				LicenseExpression        string    `json:"licenseExpression"`
				LicenseURL               string    `json:"licenseUrl"`
				Listed                   bool      `json:"listed"`
				MinClientVersion         string    `json:"minClientVersion"`
				PackageContent           string    `json:"packageContent"`
				ProjectURL               string    `json:"projectUrl"`
				Published                time.Time `json:"published"`
				RequireLicenseAcceptance bool      `json:"requireLicenseAcceptance"`
				Summary                  string    `json:"summary"`
				Tags                     []string  `json:"tags"`
				Title                    string    `json:"title"`
				Version                  string    `json:"version"`
			} `json:"catalogEntry"`
			PackageContent string `json:"packageContent"`
			Registration   string `json:"registration"`
		} `json:"items"`
		Parent string `json:"parent"`
		Lower  string `json:"lower"`
		Upper  string `json:"upper"`
	} `json:"items"`
	Context struct {
		Vocab   string `json:"@vocab"`
		Catalog string `json:"catalog"`
		Xsd     string `json:"xsd"`
		Items   struct {
			ID        string `json:"@id"`
			Container string `json:"@container"`
		} `json:"items"`
		CommitTimeStamp struct {
			ID   string `json:"@id"`
			Type string `json:"@type"`
		} `json:"commitTimeStamp"`
		CommitID struct {
			ID string `json:"@id"`
		} `json:"commitId"`
		Count struct {
			ID string `json:"@id"`
		} `json:"count"`
		Parent struct {
			ID   string `json:"@id"`
			Type string `json:"@type"`
		} `json:"parent"`
		Tags struct {
			ID        string `json:"@id"`
			Container string `json:"@container"`
		} `json:"tags"`
		Reasons struct {
			Container string `json:"@container"`
		} `json:"reasons"`
		PackageTargetFrameworks struct {
			ID        string `json:"@id"`
			Container string `json:"@container"`
		} `json:"packageTargetFrameworks"`
		DependencyGroups struct {
			ID        string `json:"@id"`
			Container string `json:"@container"`
		} `json:"dependencyGroups"`
		Dependencies struct {
			ID        string `json:"@id"`
			Container string `json:"@container"`
		} `json:"dependencies"`
		PackageContent struct {
			Type string `json:"@type"`
		} `json:"packageContent"`
		Published struct {
			Type string `json:"@type"`
		} `json:"published"`
		Registration struct {
			Type string `json:"@type"`
		} `json:"registration"`
	} `json:"@context"`
}

func (Nuget) ParseEmbeddedModules(syft *internal.Syft) (model.BuildInfo, error) {

	modules, err := Nuget.SyftToModule(Nuget{}, syft)
	if err != nil {
		return model.BuildInfo{}, err
	}

	return model.BuildInfo{
		Path:    syft.Artifacts[0].Locations[0].Path,
		Mod:     "Mod",
		Modules: Nuget.toModule(Nuget{}, modules),
	}, nil
}

func (Nuget) SyftToModule(syft *internal.Syft) ([]Module, error) {
	var result []Module
	for _, data := range syft.Artifacts {
		data.Name, _ = Nuget.createPath(Nuget{}, data.Name)
		next := Module{
			Path: data.Name,
			//subpath is maybe not needed
			SubPath: data.Version[:strings.Index(data.Version, ".")],
			Version: data.Version,
			Hash:    data.ID,
		}

		result = append(result, next)
	}

	return result, nil
}

func (Nuget) toModule(modules []Module) []model.Module {
	var result []model.Module

	for _, m := range modules {
		result = append(result, model.Module{
			Path:    m.Path,
			SubPath: m.SubPath,
			Version: m.Version,
			Hash:    m.Hash,
		})
	}

	return result
}

func (nuget Nuget) SetRepoInfo(syft *internal.Syft, info *model.BuildInfo) {
	for i := range info.Modules {
		module := &info.Modules[i]
		url := nuget.CreateAPILink(syft.Artifacts[i].Name)
		fmt.Println("[", color.Colorize(color.Yellow, "Fetch"), "] Module: ", module.Path, "from: ", url)
		pkgData, err := nuget.GetData(url)
		if err != nil {
			log.Print(err)
		}
		nuget.SetInfoToModule(module, pkgData)
	}
	fmt.Println("[", color.Colorize(color.Green, "Succ"), "] All Modules were parsed ")
}

func (Nuget) GetData(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {

		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	return body, nil
}

func (api Nuget) SetInfoToModule(module *model.Module, pkgData []byte) {
	err := json.Unmarshal(pkgData, &api)
	if err != nil {
		fmt.Println("[", color.Colorize(color.Red, "Err"), "] ", err)
	}

	//double iteration because big dependencies have multiple "sites"
	//on most dependencies there is only one api.items but not always
	for _, item := range api.Items {
		for _, data := range item.Items {
			if module.Version == data.CatalogEntry.Version {
				dep := data.CatalogEntry

				//Microsofts package descriptions start with a summary and then have a ton of
				//unimportant information. Thankfully Microsoft has a \n after the summary.
				if strings.Contains(dep.Description, "\n") {
					dep.Description = dep.Description[:strings.Index(dep.Description, "\n")]
				}

				//Some Microsoft packages that are older have the generall dot.net url and
				//not the url to the actual project. So we just construct the url ourselves.
				if dep.ProjectURL == "https://dot.net/" {
					dep.ProjectURL = fmt.Sprintf("https://www.nuget.org/packages/%s", dep.ID)
				}

				module.Info.Description = dep.Description
				module.Info.FullName = dep.Authors
				module.Info.SPDX = dep.LicenseExpression
				module.Info.Release = dep.Published
			}
		}
	}
}

func (Nuget) CreateAPILink(packageName string) string {
	packageName = strings.ToLower(packageName)
	url := fmt.Sprintf("https://api.nuget.org/v3/registration5-semver1/%s/index.json", packageName)
	return url
}

func (Nuget) createPath(name string) (string, error) {
	path := fmt.Sprintf("https://www.nuget.org/packages/%s", name)
	return path, nil
}
