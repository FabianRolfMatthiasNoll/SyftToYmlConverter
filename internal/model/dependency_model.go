package model

type SBOM struct {
	Name         string   `json:"name"`
	Languages    []string `json:"languages"`
	Dependencies []Dependencies
}

type Dependencies struct {
	ID         string   `json:"id"`
	ImportName string   `json:"importname"`
	Version    string   `json:"version"`
	Licenses   []string `json:"licenses"`
	Language   string   `json:"language"`
	//Not sure if we should keep purl for tool independency
	Purl string `json:"purl"`
	//Url      string `json:"url"`
	//TopLevel bool   `json:"toplevel"`

}
