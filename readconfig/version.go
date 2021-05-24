package readconfig

import (
	"encoding/json"
	"io/ioutil"
)

var Version VersionType

var (
	// BuildTime is a time label of the moment when the binary was built
	BuildTime string = "unset"
	// Commit is a last commit hash at the moment when the binary was built
	Commit = "unset"
	// Release is a semantic version of current build
	Release = "unset"
)

type VersionType struct {
	Minor     string `json:"Minor"`
	Major     string `json:"Major"`
	Patch     string `json:"Patch"`
	Commit    string `json:"Commit"`
	BuildTime string `json:"BuildTime"`
}

func (v *VersionType) VersionStr() string {

	return v.Major + "." + v.Minor + "." + v.Patch
}

func (v *VersionType) ReadVersionFile(versionFile string) (err error) {
	var file []byte

	if file, err = ioutil.ReadFile(versionFile); err != nil {
		return err
	}
	err = json.Unmarshal(file, &Version)
	if err != nil {
		return err
	}
	return nil
}
