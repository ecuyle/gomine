package servermanager

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/google/uuid"
)

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

// MCServer struct
type MCServer struct {
	Properties ServerProperties
	ID         string
	Name       string
	Path       string
}

// MakeWorld creates a new directory in the worlds/ directory. This new directory represents
// a new server world and will contain all necessary server files (ie. eula.txt, server.properties,
// server jarFile). After creating the new directory with the given uuid name, the appropriate
// jarFile corresponding with the provided versionID will be copied into the world and the jarFile
// will be run to instantiate required server files.
//
// The path to this new directory is returned upon successful operation.
func MakeWorld(uuid string, jarFileName string) (string, error) {
	worldPath := GetServerFilepath(uuid)
	jarFilePath := GetJarFilepath(jarFileName)

	log.Println(fmt.Sprintf("Creating world at `%v`", worldPath))
	if err := exec.Command("mkdir", "-p", worldPath).Run(); err != nil {
		return "", err
	}

	log.Println(fmt.Sprintf("Copying server jarFile from `%v` into `%v`", jarFilePath, worldPath))
	if err := exec.Command("cp", jarFilePath, worldPath).Run(); err != nil {
		return "", err
	}

	log.Println(fmt.Sprintf("Initializing server jarFile at `%v`...", worldPath))
	cmd := exec.Command("java", "-jar", jarFileName)
	cmd.Dir = worldPath
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return "", err
	}

	log.Println(fmt.Sprintf("Server jarFile successfully initialized at `%v`.", worldPath))

	return worldPath, nil
}

// MakeServer creates a server world directory for a user to later manage
func MakeServer(versionID string, serverName string, hasAcceptedEULA bool, customServerProperties map[string]interface{}) (*MCServer, error) {
	// TODO: This can all probably be cached
	version, err := GetVersionByID(versionID)

	if err != nil {
		return nil, err
	}

	versionDetails, err := GetVersionDetail(version.URL)

	if err != nil {
		return nil, err
	}

	jarFileName, err := DownloadJarFileIfNeeded(*versionDetails)

	if err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	worldPath, err := MakeWorld(id.String(), jarFileName)

	if err != nil {
		return nil, err
	}

	if err := UpdateEULA(hasAcceptedEULA, fmt.Sprintf("%v/eula.txt", worldPath)); err != nil {
		return nil, err
	}

	_, err = UpdateServerProperties(customServerProperties, fmt.Sprintf("%v/server.properties", worldPath))
	if err != nil {
		return nil, err
	}

	server := MCServer{ID: id.String(), Name: serverName, Path: worldPath}
	return &server, nil
}
