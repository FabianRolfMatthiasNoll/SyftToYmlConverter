package model

type Dependencies struct {
	Name         string         `json:"name"`
	Languages    map[string]int `json:"languages"`
	Dependencies []struct {
		Name     string   `json:"name"`
		Language string   `json:"language"`
		Version  string   `json:"version"`
		Licenses []string `json:"licenses"`
		Purl     string   `json:"purl"`
		ID       string   `json:"id"`
	}
}
