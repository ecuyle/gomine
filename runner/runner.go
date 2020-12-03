package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Version struct
type Version struct {
	ID          string `json:"id"`
	VersionType string `json:"type"`
	URL         string `json:"url"`
	Time        string `json:"time"`
	ReleaseTime string `json:"releaseTime"`
}

// LatestVersion struct
type LatestVersion struct {
	Release  string `json:"release"`
	Snapshot string `json:"snapshot"`
}

// VersionManifest struct
type VersionManifest struct {
	Latest   LatestVersion `json:"latest"`
	Versions []Version     `json:"versions"`
}

// GetVersionManifest returns the Mojang versions manifest
func GetVersionManifest() *VersionManifest {
	manifestURL := "https://launchermeta.mojang.com/mc/game/version_manifest.json"

	resp, err := http.Get(manifestURL)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	var versionManifest VersionManifest
	json.Unmarshal(body, &versionManifest)

	return &versionManifest
}

func main() {
	versionManifest := GetVersionManifest()

	log.Println(versionManifest.Latest.Release)
	log.Println(versionManifest.Versions[0].URL)
}
