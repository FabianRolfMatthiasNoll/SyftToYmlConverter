package model

type SBOM struct {
	Name         string   `json:"name"`
	Languages    []string `json:"languages"`
	Dependencies []Dependencies
}

type Dependencies struct {
	ImportName string   `json:"importname"`
	Language   string   `json:"language"`
	Version    string   `json:"version"`
	Licenses   []string `json:"licenses"`
	//Not sure if we should keep purl for tool independency
	Purl string `json:"purl"`
	//Url      string `json:"url"`
	ID string `json:"id"`
	//TopLevel bool   `json:"toplevel"`
}
