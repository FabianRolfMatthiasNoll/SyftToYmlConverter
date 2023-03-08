package api_interfaces

import (
	"fmt"
	"os/exec"
	"strings"
)

type ConanInfo struct {
	Name           string
	Version        string
	URL            string
	Homepage       string
	License        string
	Author         string
	Description    string
	Topics         string
	Generators     string
	Exports        string
	ExportsSources []string
	ShortPaths     bool
	ApplyEnv       bool
	BuildPolicy    string
	RevisionMode   string
	Settings       string
	Options        map[string][]string
	DefaultOptions map[string]string
}

func GetMetadata(packageName string, version string) ([]byte, error) {

	cmd := exec.Command("conan", "inspect", fmt.Sprintf("%s/%s@", packageName, version))
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err.Error())
		return nil, err
	}
	// conanInfo := parseConanOutput(string(output))
	return output, nil
}

func ParseConanOutput(output string) ConanInfo {
	info := ConanInfo{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "name: ") {
			info.Name = strings.TrimPrefix(line, "name: ")
		} else if strings.HasPrefix(line, "version: ") {
			info.Version = strings.TrimPrefix(line, "version: ")
		} else if strings.HasPrefix(line, "url: ") {
			info.URL = strings.TrimPrefix(line, "url: ")
		} else if strings.HasPrefix(line, "homepage: ") {
			info.Homepage = strings.TrimPrefix(line, "homepage: ")
		} else if strings.HasPrefix(line, "license: ") {
			info.License = strings.TrimPrefix(line, "license: ")
		} else if strings.HasPrefix(line, "author: ") {
			info.Author = strings.TrimPrefix(line, "author: ")
		} else if strings.HasPrefix(line, "description: ") {
			info.Description = strings.TrimPrefix(line, "description: ")
		} else if strings.HasPrefix(line, "topics: ") {
			info.Topics = strings.TrimPrefix(line, "topics: ")
		} else if strings.HasPrefix(line, "generators: ") {
			info.Generators = strings.TrimPrefix(line, "generators: ")
		} else if strings.HasPrefix(line, "exports: ") {
			info.Exports = strings.TrimPrefix(line, "exports: ")
		} else if strings.HasPrefix(line, "exports_sources: ") {
			info.ExportsSources = strings.Split(strings.TrimPrefix(line, "exports_sources: "), ", ")
		} else if strings.HasPrefix(line, "short_paths: ") {
			info.ShortPaths = strings.TrimSuffix(strings.TrimPrefix(line, "short_paths: "), "\n") == "True"
		} else if strings.HasPrefix(line, "apply_env: ") {
			info.ApplyEnv = strings.TrimSuffix(strings.TrimPrefix(line, "apply_env: "), "\n") == "True"
		} else if strings.HasPrefix(line, "build_policy: ") {
			info.BuildPolicy = strings.TrimPrefix(line, "build_policy: ")
		} else if strings.HasPrefix(line, "revision_mode: ") {
			info.RevisionMode = strings.TrimPrefix(line, "revision_mode: ")
		} else if strings.HasPrefix(line, "settings: ") {
			info.Settings = strings.TrimPrefix(line, "settings: ")
		} else if strings.HasPrefix(line, "options:") {
			// Parse options section
			info.Options = make(map[string][]string)
			for i := 1; i < len(lines); i++ {
				line := strings.TrimSpace(lines[i])
				if line == "" {
					break
				}
				option := strings.Split(line, ":")
				name := strings.TrimSpace(option[0])
				values := strings.TrimSpace(option[1])
				info.Options[name] = strings.Split(values[1:len(values)-1], ", ")
			}
		} else if strings.HasPrefix(line, "default_options:") {
			// Parse default options section
			info.DefaultOptions = make(map[string]string)
			for i := 1; i < len(lines); i++ {
				line := strings.TrimSpace(lines[i])
				if line == "" {
					break
				}
				option := strings.Split(line, ":")
				name := strings.TrimSpace(option[0])
				value := strings.TrimSpace(option[1])
				info.DefaultOptions[name] = value
			}
		}
	}
	return info
}
