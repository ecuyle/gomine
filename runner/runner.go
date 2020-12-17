package main

import (
	"encoding/json"
	"errors"
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

// VersionDownload struc
type VersionDownload struct {
	Sha1 string `json:"sha1"`
	Size int    `json:"size"`
	URL  string `json:"url"`
}

// VersionDetail struct
type VersionDetail struct {
	Downloads map[string]VersionDownload `json:"downloads"`
	ID        string                     `json:"id"`
}

// GetVersionByID returns the Mojang version object given a version id
func GetVersionByID(versionID string) (*Version, error) {
	versionManifest, err := GetVersionManifest()

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(versionManifest.Versions); i++ {
		var version = versionManifest.Versions[i]

		if versionID == version.ID {
			return &version, nil
		}
	}

	return nil, errors.New("Could not find version with id: " + versionID)
}

// GetVersionManifest returns the Mojang versions manifest
func GetVersionManifest() (*VersionManifest, error) {
	manifestUrl := "https://launchermeta.mojang.com/mc/game/version_manifest.json"

	resp, err := http.Get(manifestUrl)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var versionManifest VersionManifest
	err = json.Unmarshal(body, &versionManifest)

	if err != nil {
		return nil, err
	}

	return &versionManifest, nil
}

func main() {
	versionManifest, err := GetVersionManifest()

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Latest Release:", versionManifest.Latest.Release)

	versionID := "1.16.4"

	version, err := GetVersionByID(versionID)

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(version)
	log.Println(version.URL)
}
