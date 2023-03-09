package model

type SBOM struct {
	Name         string
	Languages    []string
	Dependencies []Dependencies
}

type Dependencies struct {
	ID         string
	ImportName string
	Version    string
	Licenses   []string
	Language   string
	//Not sure if we should keep purl for tool independency
	Purl string
	//Url      string
	//TopLevel bool

}
