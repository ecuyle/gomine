package servermanager

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

// ServerProperties struct
type ServerProperties struct {
	AllowFlight                    bool   `alias:"allow-flight" json:"allow-flight" properties:"allow-flight,default=false"`                                                              // false
	AllowNether                    bool   `alias:"allow-nether" json:"allow-nether" properties:"allow-nether,default=true"`                                                               // true
	BroadcastConsoleToOps          bool   `alias:"broadcast-console-to-ops" json:"broadcast-console-to-ops" properties:"broadcast-console-to-ops,default=true"`                           // true
	BroadcastRconToOps             bool   `alias:"broadcast-rcon-to-ops" json:"broadcast-rcon-to-ops" properties:"broadcast-rcon-to-ops,default=true"`                                    // true
	ConPort                        uint16 `alias:"con.port" json:"con.port" properties:"con.port,default=25575"`                                                                          // 25575
	Difficulty                     string `alias:"difficulty" json:"difficulty" properties:"difficulty,default=easy"`                                                                     // easy
	EnableCommandBlock             bool   `alias:"enable-command-block" json:"enable-command-block" properties:"enable-command-block,default=false"`                                      // false
	EnableJMXMonitoring            bool   `alias:"enable-jmx-monitoring" json:"enable-jmx-monitoring" properties:"enable-jmx-monitoring,default=false"`                                   // false
	EnableQuery                    bool   `alias:"enable-query" json:"enable-query" properties:"enable-query,default=false"`                                                              // false
	EnableRcon                     bool   `alias:"enable-rcon" json:"enable-rcon" properties:"enable-rcon,default=false"`                                                                 // false
	EnableStatus                   bool   `alias:"enable-status" json:"enable-status" properties:"enable-status,default=true"`                                                            // true
	EnforceWhitelist               bool   `alias:"enforce-whitelist" json:"enforce-whitelist" properties:"enforce-whitelist,default=false"`                                               // false
	EntityBroadcastRangePercentage int    `alias:"entity-broadcast-range-percentage" json:"entity-broadcast-range-percentage" properties:"entity-broadcast-range-percentage,default=100"` // 100
	ForceGamemode                  bool   `alias:"force-gamemode" json:"force-gamemode" properties:"force-gamemode,default=false"`                                                        // false
	FunctionPermissionLevel        int    `alias:"function-permission-level" json:"function-permission-level" properties:"function-permission-level,default=2"`                           // 2
	Gamemode                       string `alias:"gamemode" json:"gamemode" properties:"gamemode,default=survival"`                                                                       // survival
	GenerateStructures             bool   `alias:"generate-structures" json:"generate-structures" properties:"generate-structures,default=true"`                                          // true
	GeneratorSettings              string `alias:"generator-settings" json:"generator-settings,omitempty" properties:"generator-settings,default="`                                       //
	Hardcore                       bool   `alias:"hardcore" json:"hardcore" properties:"hardcore,default=false"`                                                                          // false
	LevelName                      string `alias:"level-name" json:"level-name" properties:"level-name,default=world"`                                                                    // world
	LevelSeed                      string `alias:"level-seed" json:"level-seed" properties:"level-seed,default="`                                                                         //
	LevelType                      string `alias:"level-type" json:"level-type" properties:"level-type,default=default"`                                                                  // default
	MaxBuildHeight                 int    `alias:"max-build-height" json:"max-build-height" properties:"max-build-height,default=256"`                                                    // 256
	MaxPlayers                     int    `alias:"max-players" json:"max-players" properties:"max-players,default=20"`                                                                    // 20
	MaxTickTime                    int32  `alias:"max-tick-time" json:"max-tick-time" properties:"max-tick-time,default=60000"`                                                           // 60000
	MaxWorldSize                   int64  `alias:"max-world-size" json:"max-world-size" properties:"max-world-size,default=29999984"`                                                     // 29999984
	Motd                           string `alias:"motd" json:"motd" properties:"motd,default=A Minecraft Server"`                                                                         // A Minecraft Server
	NetworkCompressionThreshold    int    `alias:"network-compression-threshold" json:"network-compression-threshold" properties:"network-compression-threshold,default=256"`             // 256
	OnlineMode                     bool   `alias:"online-mode" json:"online-mode" properties:"online-mode,default=true"`                                                                  // true
	OpPermissionLevel              int    `alias:"op-permission-level" json:"op-permission-level" properties:"op-permission-level,default=4"`                                             // 4
	PVP                            bool   `alias:"pvp" json:"pvp" properties:"pvp,default=true"`                                                                                          // true
	PlayerIdleTimeout              int    `alias:"player-idle-timeout" json:"player-idle-timeout" properties:"player-idle-timeout,default=0"`                                             // 0
	PreventProxyConnections        bool   `alias:"prevent-proxy-connections" json:"prevent-proxy-connections" properties:"prevent-proxy-connections,default=falsefalse"`                  // false
	QueryPort                      uint16 `alias:"query.port" json:"query.port" properties:"query.port,default=25565"`                                                                    // 25565
	RateLimit                      int    `alias:"rate-limit" json:"rate-limit" properties:"rate-limit,default=0"`                                                                        // 0
	RconPassword                   string `alias:"rcon.password" json:"rcon.password,omitempty" properties:"rcon.password,default="`                                                      //
	ResourcePack                   string `alias:"resource-pack" json:"resource-pack,omitempty" properties:"resource-pack,default="`                                                      //
	ResourcePackSha1               string `alias:"resource-pack-sha1" json:"resource-pack-sha1,omitempty" properties:"resource-pack-sha1,default="`                                       //
	ServerIP                       string `alias:"server-ip" json:"server-ip,omitempty" properties:"server-ip,default="`                                                                  //
	ServerPort                     uint16 `alias:"server-port" json:"server-port" properties:"server-port,default=25565"`                                                                 // 25565
	SnooperEnabled                 bool   `alias:"snooper-enabled" json:"snooper-enabled" properties:"snooper-enabled,default=true"`                                                      // true
	SpawnAnimals                   bool   `alias:"spawn-animals" json:"spawn-animals" properties:"spawn-animals,default=true"`                                                            // true
	SpawnMonsters                  bool   `alias:"spawn-monsters" json:"spawn-monsters" properties:"spawn-monsters,default=true"`                                                         // true
	SpawnNpcs                      bool   `alias:"spawn-npcs" json:"spawn-npcs" properties:"spawn-npcs,default=true"`                                                                     // true
	SpawnProtection                int    `alias:"spawn-protection" json:"spawn-protection" properties:"spawn-protection,default=16"`                                                     // 16
	SyncChunkWrites                bool   `alias:"sync-chunk-writes" json:"sync-chunk-writes" properties:"sync-chunk-writes,default=true"`                                                // true
	UseNativeTransport             bool   `alias:"use-native-transport" json:"use-native-transport" properties:"use-native-transport,default=true"`                                       // true
	ViewDistance                   int    `alias:"view-distance" json:"view-distance" properties:"view-distance,default=10"`                                                              // 10
	WhiteList                      bool   `alias:"white-list" json:"white-list" properties:"white-list,default=false"`                                                                    // false
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

	defer resp.Body.Close()
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
func DownloadFile(url string, filepath string) error {
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

const DATA_PATH_PREFIX = "data/"

// GetJarFilepath returns the standard location for downloaded server jarFiles
func GetJarFilepath(jarFileName string) string {
	return fmt.Sprintf("%vjarFiles/%v", DATA_PATH_PREFIX, jarFileName)
}

// GetServerFilePath gets the full path to a server from the worlds/
// directory starting at the project root
func GetServerFilepath(serverID string) string {
	return fmt.Sprintf("%vworlds/%v", DATA_PATH_PREFIX, serverID)
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

	if err := DownloadFile(jarFileURL, jarFilePath); err != nil {
		return "", err
	}

	return jarFileName, nil
}

func GetEULAFilepath(worldpath string) string {
	return fmt.Sprintf("%v/eula.txt", worldpath)
}

func IsEulaAccepted(worldpath string) bool {
	data, err := ioutil.ReadFile(GetEULAFilepath(worldpath))

	if err != nil {
		return false
	}

	trueInBytes := []byte(strconv.FormatBool(true))

	return bytes.Contains(data, trueInBytes)
}

// UpdateEULA updates the eula.txt for a server with the provided value
func UpdateEULA(value bool, worldpath string) error {
	filepath := GetEULAFilepath(worldpath)
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

func GetServerPropertiesFilepath(worldpath string) string {
	return fmt.Sprintf("%v/server.properties", worldpath)
}

// GetServerProperties gets the current server properties for a given server
func GetServerProperties(worldpath string) (*properties.Properties, error) {
	result, err := properties.LoadFile(GetServerPropertiesFilepath(worldpath), properties.UTF8)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// WriteServerProperties updates a server.properties file with an updated properties.Properties struct
func WriteServerProperties(worldpath string, serverProperties *properties.Properties) error {
	serverPropertiesFile, err := os.Create(GetServerPropertiesFilepath(worldpath))
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
func UpdateServerProperties(customServerProperties map[string]interface{}, worldpath string) (*ServerProperties, error) {
	currentServerProperties, err := GetServerProperties(worldpath)

	if err != nil {
		return nil, err
	}

	for key, value := range customServerProperties {
		if err := currentServerProperties.SetValue(key, value); err != nil {
			return nil, err
		}
	}

	if err := WriteServerProperties(worldpath, currentServerProperties); err != nil {
		return nil, err
	}

	log.Println(fmt.Sprintf("`%v` updated with new values: %v", GetServerPropertiesFilepath(worldpath), customServerProperties))

	updatedServerProperties := ServerProperties{}
	if err := currentServerProperties.Decode(&updatedServerProperties); err != nil {
		return nil, err
	}

	return &updatedServerProperties, nil
}
