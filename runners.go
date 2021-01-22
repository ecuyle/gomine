package main

import (
	"bufio"
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

	"github.com/magiconair/properties"
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

// GetJarFilepath returns the standard location for downloaded server jarFiles
func GetJarFilepath(jarFileName string) string {
	return fmt.Sprintf("jarFiles/%v", jarFileName)
}

// GetServerFilePath gets the full path to a server from the worlds/
// directory starting at the project root
func GetServerFilepath(serverID string) string {
	return fmt.Sprintf("worlds/%v", serverID)
}

// DownloadJarFileIfNeeded download a jarFile if the desired jarFile has not already been downloaded
func DownloadJarFileIfNeeded(versionDetail VersionDetail) (string, error) {
	jarFileName := fmt.Sprintf("%v.jar", versionDetail.ID)
	jarFilePath := GetJarFilepath(jarFileName)

	if _, err := os.Stat(jarFilePath); err == nil {
		log.Println(fmt.Sprintf("jarFile `%v` found. Skipping download.", jarFilePath))
		return jarFileName, nil
	}

	jarFileURL := versionDetail.Downloads.Server.URL
	log.Println(fmt.Sprintf("Downloading jarFile from `%v` into `%v`", jarFileURL, jarFilePath))

	if err := exec.Command("mkdir", "-p", "jarFiles").Run(); err != nil {
		return "", err
	}

	if err := DownloadFile(jarFilePath, jarFileURL); err != nil {
		return "", err
	}

	return jarFileName, nil
}

// UpdateEULA updates the eula.txt for a server with the provided value
func UpdateEULA(value bool, filepath string) error {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return err
	}

	oldValueInBytes := []byte(strconv.FormatBool(!value))

	if !bytes.Contains(data, oldValueInBytes) {
		log.Println(fmt.Sprintf("`%v` already is set to value `%v`. Skipping EULA update.", filepath, value))
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

	log.Println(fmt.Sprintf("`%v` updated with new value `%v`", filepath, value))
	return nil
}

// GetServerProperties gets the current server properties for a given server
func GetServerProperties(filepath string) *properties.Properties {
	return properties.MustLoadFile(filepath, properties.UTF8)
}

// WriteServerProperties updates a server.properties file with an updated properties.Properties struct
func WriteServerProperties(filepath string, serverProperties *properties.Properties) error {
	serverPropertiesFile, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer serverPropertiesFile.Close()

	writer := bufio.NewWriter(serverPropertiesFile)
	if _, err := serverProperties.Write(writer, properties.UTF8); err != nil {
		return err
	}
	writer.Flush()

	return nil
}

// UpdateServerProperties updates the server.properties file for a given server
func UpdateServerProperties(customServerProperties map[string]interface{}, filepath string) (*ServerProperties, error) {
	currentServerProperties := GetServerProperties(filepath)

	for key, value := range customServerProperties {
		if err := currentServerProperties.SetValue(key, value); err != nil {
			return nil, err
		}
	}

	if err := WriteServerProperties(filepath, currentServerProperties); err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("`%v` updated with new values: %v", filepath, customServerProperties))

	updatedServerProperties := ServerProperties{}
	if err := currentServerProperties.Decode(&updatedServerProperties); err != nil {
		return nil, err
	}

	return &updatedServerProperties, nil
}
