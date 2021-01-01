package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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
	EnableJMXMonitoring            bool   `json:"enable-jmx-monitoring"`             // false
	ConPort                        int    `json:"con.port"`                          // 25575
	LevelSeed                      string `json:"level-seed"`                        //
	Gamemode                       string `json:"gamemode"`                          // survival
	EnableCommandBlock             bool   `json:"enable-command-block"`              // false
	EnableQuery                    bool   `json:"enable-query"`                      // false
	GeneratorSettings              string `json:"generator-settings,omitempty"`      //
	LevelName                      string `json:"level-name"`                        // world
	Motd                           string `json:"motd"`                              // A Minecraft Server
	QueryPort                      int    `json:"query.port"`                        // 25565
	PVP                            bool   `json:"pvp"`                               // true
	GenerateStructures             bool   `json:"generate-structures"`               // true
	Difficulty                     string `json:"difficulty"`                        // easy
	NetworkCompressionThreshold    int    `json:"network-compression-threshold"`     // 256
	MaxTickTime                    int    `json:"max-tick-time"`                     // 60000
	MaxPlayers                     int    `json:"max-players"`                       // 20
	UseNativeTransport             bool   `json:"use-native-transport"`              // true
	OnlineMode                     bool   `json:"online-mode"`                       // true
	EnableStatus                   bool   `json:"enable-status"`                     // true
	AllowFlight                    bool   `json:"allow-flight"`                      // false
	BroadcastRconToOps             bool   `json:"broadcast-rcon-to-ops"`             // true
	ViewDistance                   int    `json:"view-distance"`                     // 10
	MaxBuildHeight                 int    `json:"max-build-height"`                  // 256
	ServerIP                       string `json:"server-ip,omitempty"`               //
	AllowNether                    bool   `json:"allow-nether"`                      // true
	ServerPort                     int    `json:"server-port"`                       // 25565
	EnableRcon                     bool   `json:"enable-rcon"`                       // false
	SyncChunkWrites                bool   `json:"sync-chunk-writes"`                 // true
	OpPermissionLevel              int    `json:"op-permission-level"`               // 4
	PreventProxyConnections        bool   `json:"prevent-proxy-connections"`         // false
	ResourcePack                   string `json:"resource-pack,omitempty"`           //
	EntityBroadcastRangePercentage int    `json:"entity-broadcast-range-percentage"` // 100
	RconPassword                   string `json:"rcon.password,omitempty"`           //
	PlayerIdleTimeout              int    `json:"player-idle-timeout"`               // 0
	ForceGamemode                  bool   `json:"force-gamemode"`                    // false
	RateLimit                      int    `json:"rate-limit"`                        // 0
	Hardcore                       bool   `json:"hardcore"`                          // false
	WhiteList                      bool   `json:"white-list"`                        // false
	BroadcastConsoleToOps          bool   `json:"broadcast-console-to-ops"`          // true
	SpawnNpcs                      bool   `json:"spawn-npcs"`                        // true
	SpawnAnimals                   bool   `json:"spawn-animals"`                     // true
	SnooperEnabled                 bool   `json:"snooper-enabled"`                   // true
	FunctionPermissionLevel        int    `json:"function-permission-level"`         // 2
	LevelType                      string `json:"level-type"`                        // default
	SpawnMonsters                  bool   `json:"spawn-monsters"`                    // true
	EnforceWhitelist               bool   `json:"enforce-whitelist"`                 // false
	ResourcePackSha1               string `json:"resource-pack-sha1,omitempty"`      //
	SpawnProtection                int    `json:"spawn-protection"`                  // 16
	MaxWorldSize                   int    `json:"max-world-size"`                    // 29999984
}

// MCServer struct
type MCServer struct {
	Properties ServerProperties
	ID         string
	Name       string
	Path       string
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

	exec.Command("mkdir", "-p", "jarfiles").Run()
	if err := DownloadFile(jarfilePath, jarfileURL); err != nil {
		return "", err
	}

	return jarfileName, nil
}

// CreateWorld creates a new directory in the worlds/ directory. This new directory represents
// a new server world and will contain all necessary server files (ie. eula.txt, server.properties,
// server jarfile). After creating the new directory with the given uuid name, the appropriate
// jarfile corresponding with the provided versionID will be copied into the world and the jarfile
// will be run to instantiate required server files.
//
// The path to this new directory is returned upon successful operation.
func CreateWorld(uuid string, jarfileName string) (string, error) {
	worldPath := fmt.Sprintf("worlds/%v", uuid)
	jarfilePath := GetJarfilePath(jarfileName)

	log.Println(fmt.Sprintf("Creating world at %v", worldPath))
	if err := exec.Command("mkdir", "-p", worldPath).Run(); err != nil {
		return "", err
	}

	log.Println(fmt.Sprintf("Copying server jarfile from %v into %v", jarfilePath, worldPath))
	if err := exec.Command("cp", jarfilePath, worldPath).Run(); err != nil {
		return "", err
	}

	log.Println("Initializing server jarfile...")
	cmd := exec.Command("java", "-jar", jarfileName)
	cmd.Dir = worldPath
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("Success.")

	return worldPath, nil
}

// CreateServer creates a server world directory for a user to later manage
func CreateServer(versionID string, serverName string, isEulaAccepted bool, serverProperties ServerProperties) (*MCServer, error) {
	// TODO: This can all probably be cached
	version, err := GetVersionByID(versionID)

	if err != nil {
		return nil, err
	}

	versionDetails, err := GetVersionDetail(version.URL)

	if err != nil {
		return nil, err
	}

	jarfileName, err := DownloadJarfileIfNeeded(*versionDetails)

	if err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	worldPath, err := CreateWorld(id.String(), jarfileName)

	if err != nil {
		return nil, err
	}

	// TODO: Update eula.txt and server.properties
	server := MCServer{ID: id.String(), Name: serverName, Path: worldPath}
	return &server, nil
}

func main() {
	var serverProperties ServerProperties
	_, err := CreateServer("1.16.4", "TestServer", true, serverProperties)

	if err != nil {
		log.Fatalln(err)
	}
}
