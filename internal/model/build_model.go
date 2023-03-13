package model

import (
	"fmt"
	"time"
)

type BuildInfo struct {
	Path    string
	Mod     string
	Modules []Module
}

type Module struct {
	Name    string
	Path    string
	SubPath string
	Version string
	Hash    string
	Parents []string
	Info    RepoInfo
}

func (m Module) String() string {
	if m.SubPath != "" {
		return fmt.Sprintf("%s/%s@%s", m.Path, m.SubPath, m.Version)
	}

	return fmt.Sprintf("%s@%s", m.Path, m.Version)
}

type RepoInfo struct {
	FullName    string
	Description string
	SPDX        string
	Release     time.Time
}
