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

	"github.com/TwiN/go-color"
	"github.com/gammazero/workerpool"
)

type Author struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type NPM struct {
	Name          string      `json:"name"`
	Version       string      `json:"version"`
	Description   string      `json:"description"`
	Author        interface{} `json:"author"`
	Homepage      string      `json:"homepage"`
	License       string      `json:"license"`
	PublishConfig struct {
		Access string `json:"access"`
	} `json:"publishConfig"`
	Main            string            `json:"main"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	ID              string            `json:"_id"`
	Dist            struct {
		Shasum       string `json:"shasum"`
		Integrity    string `json:"integrity"`
		Tarball      string `json:"tarball"`
		FileCount    int    `json:"fileCount"`
		UnpackedSize int    `json:"unpackedSize"`
		NpmSignature string `json:"npm-signature"`
		Signatures   []struct {
			Keyid string `json:"keyid"`
			Sig   string `json:"sig"`
		} `json:"signatures"`
	} `json:"dist"`
	NpmUser struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"_npmUser"`
	Directories struct {
	} `json:"directories"`
	Maintainers []struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"maintainers"`
}

func (npm NPM) ParseEmbeddedModules(syft *internal.Syft) (model.BuildInfo, error) {

	modules, err := npm.SyftToModule(syft)
	if err != nil {
		return model.BuildInfo{}, err
	}

	return model.BuildInfo{
		Path:    syft.Artifacts[0].Locations[0].Path,
		Mod:     "Mod",
		Modules: npm.toModule(modules),
	}, nil
}

func (npm NPM) SyftToModule(syft *internal.Syft) ([]Module, error) {
	var result []Module
	for _, data := range syft.Artifacts {
		data.Name, _ = npm.createPath(data.Name)
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

func (NPM) toModule(modules []Module) []model.Module {
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

func (npm NPM) SetRepoInfo(syft *internal.Syft, info *model.BuildInfo) {
	wp := workerpool.New(10)

	for i := range info.Modules {
		module := &info.Modules[i]
		wp.Submit(func() {
			pkgName := npm.getNameFromPath(module.Path)
			url := npm.CreateAPILink(pkgName, module.Version)
			fmt.Println("[", color.Colorize(color.Yellow, "Fetch"), "] Module:", pkgName, "from:", url)
			pkgData, err := npm.GetData(url)
			if err != nil {
				log.Print(err)
			}
			npm.SetInfoToModule(module, pkgData)
		})
	}
	wp.StopWait()
	fmt.Println("[", color.Colorize(color.Green, "Succ"), "] All Modules were parsed ")
}

func (NPM) GetData(url string) ([]byte, error) {
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

func (api NPM) SetInfoToModule(module *model.Module, pkgData []byte) error {
	err := json.Unmarshal(pkgData, &api)
	if err != nil {
		fmt.Println("[", color.Colorize(color.Red, "Err"), "] ", err)
		return err
	}
	module.Info.Description = api.Description

	// Sometimes Author is a struct and sometimes a string => Created Owner as author string
	// and Author as author struct
	if authorString, ok := api.Author.(string); ok {
		module.Info.FullName = authorString

	}
	if authorStruct, ok := api.Author.(Author); ok {
		module.Info.FullName = authorStruct.Name
	}

	module.Info.SPDX = api.License
	//cant extract time of release of version specific package
	//module.Info.Release =
	return nil
}

func (NPM) CreateAPILink(packageName, version string) string {
	packageName = strings.ToLower(packageName)
	url := fmt.Sprintf("https://registry.npmjs.org/%s/%s", packageName, version)
	return url
}

func (NPM) createPath(name string) (string, error) {
	path := fmt.Sprintf("https://www.npmjs.com/package/%s", name)
	return path, nil
}

func (NPM) getNameFromPath(path string) string {
	path = path[30:]
	return path
}

func (npm NPM) SetParents(model *model.BuildInfo) {
	for i := range model.Modules {
		module := &model.Modules[i]
		pkgNameParent := npm.getNameFromPath(module.Path)
		fmt.Println("[", color.Colorize(color.Yellow, "Info"), "] Scanning Dependencies of:", pkgNameParent)
		url := npm.CreateAPILink(pkgNameParent, module.Version)
		pkgData, err := npm.GetData(url)
		if err != nil {
			log.Print(err)
		}
		err = json.Unmarshal(pkgData, &npm)
		if err != nil {
			fmt.Println("[", color.Colorize(color.Red, "Err"), "] ", err)
		}
		for p := range npm.Dependencies {
			fmt.Println("[", color.Colorize(color.Yellow, "Info"), "] Found Dependency:", p)
			for r := range model.Modules {
				module := &model.Modules[r]
				pkgName := npm.getNameFromPath(module.Path)
				if strings.Contains(p, pkgName) {
					if !contains(module.Parents, pkgName) {
						//TODO: functions sets itself as parent!!!
						fmt.Println("[", color.Colorize(color.Green, "Set"), "] Set Parent Dependency to:", pkgName)
						module.Parents = append(module.Parents, pkgNameParent)
					}
				}
			}
		}
		// for p := range npm.DevDependencies {
		// 	for r := range model.Modules {
		// 		module := &model.Modules[r]
		// 		pkgName := npm.getNameFromPath(module.Path)
		// 		if strings.Contains(p, pkgName) {
		// 			module.Parents = append(module.Parents, pkgName)
		// 		}
		// 	}
		// }
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
