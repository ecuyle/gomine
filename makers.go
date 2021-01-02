package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/google/uuid"
)

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

// MakeWorld creates a new directory in the worlds/ directory. This new directory represents
// a new server world and will contain all necessary server files (ie. eula.txt, server.properties,
// server jarfile). After creating the new directory with the given uuid name, the appropriate
// jarfile corresponding with the provided versionID will be copied into the world and the jarfile
// will be run to instantiate required server files.
//
// The path to this new directory is returned upon successful operation.
func MakeWorld(uuid string, jarfileName string) (string, error) {
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

// MakeServer creates a server world directory for a user to later manage
func MakeServer(versionID string, serverName string, hasAcceptedEULA bool, serverProperties ServerProperties) (*MCServer, error) {
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

	worldPath, err := MakeWorld(id.String(), jarfileName)

	if err != nil {
		return nil, err
	}

	if err := UpdateEULA(hasAcceptedEULA, fmt.Sprintf("%v/eula.txt", worldPath)); err != nil {
		return nil, err
	}

	// TODO: Update eula.txt and server.properties
	server := MCServer{ID: id.String(), Name: serverName, Path: worldPath}
	return &server, nil
}
