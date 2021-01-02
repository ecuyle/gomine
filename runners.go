package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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

// VersionDownloads struct
type VersionDownloads struct {
	Server VersionDownload `json:"server"`
}

// VersionDetail struct
type VersionDetail struct {
	Downloads VersionDownloads `json:"downloads"`
	ID        string           `json:"id"`
}

// GetVersionDetail returns the details of a given Mojang version object
func GetVersionDetail(versionURL string) (*VersionDetail, error) {
	resp, err := http.Get(versionURL)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var versionDetail VersionDetail
	err = json.Unmarshal(body, &versionDetail)

	if err != nil {
		return nil, err
	}

	return &versionDetail, nil
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
	manifestURL := "https://launchermeta.mojang.com/mc/game/version_manifest.json"

	resp, err := http.Get(manifestURL)

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

// DownloadFile downloads a file to a filepath given a url
func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	out, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	return err
}

// GetJarfilePath returns the standard location for downloaded server jarfiles
func GetJarfilePath(jarfileName string) string {
	return fmt.Sprintf("jarfiles/%v", jarfileName)
}

// DownloadJarfileIfNeeded download a jarfile if the desired jarfile has not already been downloaded
func DownloadJarfileIfNeeded(versionDetail VersionDetail) (string, error) {
	jarfileName := fmt.Sprintf("%v.jar", versionDetail.ID)
	jarfilePath := GetJarfilePath(jarfileName)

	if _, err := os.Stat(jarfilePath); err == nil {
		log.Println(fmt.Sprintf("Jarfile %v found. Skipping download.", jarfilePath))
		return jarfileName, nil
	}

	jarfileURL := versionDetail.Downloads.Server.URL
	log.Println(fmt.Sprintf("Downloading jarfile from %v into %v", jarfileURL, jarfilePath))

	if err := exec.Command("mkdir", "-p", "jarfiles").Run(); err != nil {
		return "", err
	}

	if err := DownloadFile(jarfilePath, jarfileURL); err != nil {
		return "", err
	}

	return jarfileName, nil
}

// UpdateEULA updates the eula.txt for a server with the provided value
func UpdateEULA(value bool, filepath string) error {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	oldValueInBytes := []byte(strconv.FormatBool(!value))

	if !bytes.Contains(data, oldValueInBytes) {
		log.Println(fmt.Sprintf("%v already is set to value `%v`. Skipping EULA update.", filepath, value))
		return nil
	}

	valueInBytes := []byte(strconv.FormatBool(value))
	newData := bytes.Replace(data, oldValueInBytes, valueInBytes, 1)

	file, err := os.OpenFile(filepath, os.O_RDWR, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.WriteAt(newData, 0); err != nil {
		return err
	}

	log.Println(fmt.Sprintf("%v updated with new value `%v`", filepath, value))
	return nil
}
