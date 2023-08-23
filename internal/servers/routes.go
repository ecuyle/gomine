package servers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os/exec"

	_ "github.com/mattn/go-sqlite3"

	httputils "github.com/ecuyle/gomine/internal/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/magiconair/properties"
)

func GetDefaults(context *gin.Context) {
	var defaultProperties ServerProperties
	var p properties.Properties

	p.Decode(&defaultProperties)
	context.IndentedJSON(http.StatusOK, &defaultProperties)
}

// makeWorld creates a new directory in the worlds/ directory. This new directory represents
// a new server world and will contain all necessary server files (ie. eula.txt, server.properties,
// server jarFile). After creating the new directory with the given uuid name, the appropriate
// jarFile corresponding with the provided versionID will be copied into the world and the jarFile
// will be run to instantiate required server files.
//
// The path to this new directory is returned upon successful operation.
func makeWorld(uuid string, jarFileName string) (string, error) {
	worldPath := GetServerFilepath(uuid)
	jarFilePath := GetJarFilepath(jarFileName)

	log.Printf("Creating world at `%v`", worldPath)
	if err := exec.Command("mkdir", "-p", worldPath).Run(); err != nil {
		return "", err
	}

	log.Printf("Copying server jarFile from `%v` into `%v`", jarFilePath, worldPath)
	if err := exec.Command("cp", jarFilePath, worldPath).Run(); err != nil {
		return "", err
	}

	log.Printf("Initializing server jarFile at `%v`...", worldPath)
	cmd := exec.Command("java", "-jar", jarFileName)
	cmd.Dir = worldPath
	if err := cmd.Run(); err != nil {
		log.Println(err)
		return "", err
	}

	log.Printf("Server jarFile successfully initialized at `%v`.", worldPath)

	return worldPath, nil
}

type ServerOptions struct {
	Name           string                 `json:"name"`
	UserID         string                 `json:"userId"`
	Runtime        string                 `json:"runtime"`
	IsEulaAccepted bool                   `json:"isEulaAccepted"`
	Config         map[string]interface{} `json:"config"`
}

type MCServerLite struct {
	ID      string
	Name    string
	PID     int
	Path    string
	Runtime string
	Status  bool
	UserID  string
}

// MCServer struct
type MCServer struct {
	ID             string
	IsEulaAccepted bool
	Name           string
	PID            int
	Path           string
	Properties     ServerProperties
	Runtime        string
	Status         bool
	UserID         string
}

// makeServer creates a server world directory for a user to later manage
func makeServer(options *ServerOptions) (*MCServer, error) {
	runtime := options.Runtime

	// TODO: This can all probably be cached
	version, err := GetVersionByID(runtime)

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

	worldPath, err := makeWorld(id.String(), jarFileName)

	if err != nil {
		return nil, err
	}

	isEulaAccepted := options.IsEulaAccepted

	if err := UpdateEULA(isEulaAccepted, worldPath); err != nil {
		return nil, err
	}

	updatedServerProperties, err := UpdateServerProperties(options.Config, worldPath)
	if err != nil {
		return nil, err
	}

	server := MCServer{
		ID:             id.String(),
		IsEulaAccepted: isEulaAccepted,
		Name:           options.Name,
		PID:            -1,
		Path:           worldPath,
		Properties:     *updatedServerProperties,
		Runtime:        runtime,
		Status:         false,
		UserID:         options.UserID,
	}

	return &server, nil
}

func insertServerRecord(server *MCServer) error {
	db, err := sql.Open("sqlite3", "./gomine.db")

	if err != nil {
		return err
	}

	defer db.Close()
	transaction, err := db.Begin()

	if err != nil {
		return err
	}

	statement, err := transaction.Prepare("insert into servers(id, name, runtime, path, pid, status, user_id) values(?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(server.ID, server.Name, server.Runtime, server.Path, server.PID, server.Status, server.UserID)

	if err != nil {
		return err
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}

func PostServer(context *gin.Context) {
	var options ServerOptions

	if err := context.BindJSON(&options); err != nil {
		log.Println(err)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	server, err := makeServer(&options)

	if err != nil {
		httputils.RespondWithInternalServerError(context, err)
		return
	}

	err = insertServerRecord(server)

	if err != nil {
		httputils.RespondWithInternalServerError(context, err)
		return
	}

	httputils.RespondWithStatusCreated(context, server)
}

type UpdatedServerProperties struct {
	ServerID         string                 `json:"serverId"`
	ServerProperties map[string]interface{} `json:"serverProperties"`
}

func updateServerWorld(serverId string, properties map[string]interface{}) (*ServerProperties, error) {
	filepath := GetServerFilepath(serverId)
	updatedProperties, err := UpdateServerProperties(properties, filepath)

	if err != nil {
		return nil, err
	}

	return updatedProperties, nil
}

func PutServerProperties(context *gin.Context) {
	var options UpdatedServerProperties

	if err := context.BindJSON(&options); err != nil {
		log.Println(err)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	updatedProperties, err := updateServerWorld(options.ServerID, options.ServerProperties)

	if err != nil {
		httputils.RespondWithInternalServerError(context, err)
		return
	}

	httputils.RespondWithStatusCreated(context, updatedProperties)
}

func selectServerRecordById(id string) (*MCServer, error) {
	db, err := sql.Open("sqlite3", "./gomine.db")

	if err != nil {
		return nil, err
	}

	defer db.Close()

	statement, err := db.Prepare("select name, runtime, path, pid, status, user_id from servers where id=?")

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	server := MCServer{ID: id}
	err = statement.QueryRow(id).Scan(&server.Name, &server.Runtime, &server.Path, &server.PID, &server.Status, &server.UserID)

	if err != nil {
		return nil, err
	}

	return &server, nil
}

func populateServerWithProperties(server *MCServer) error {
	properties := ServerProperties{}
	serverPropertiesFromDisk, err := GetServerProperties(server.Path)

	if err != nil {
		return err
	}

	if err := serverPropertiesFromDisk.Decode(&properties); err != nil {
		return err
	}

	server.Properties = properties

	return nil
}

func populateServerWithEulaAcceptanceStatus(server *MCServer) error {
	server.IsEulaAccepted = IsEulaAccepted(server.Path)
	return nil
}

func GetServerDetails(context *gin.Context) {
	serverId := context.Query("s")

	if serverId == "" {
		httputils.RespondWithNotFound(context, errors.New("GetServerDetails: No server id provided."))
		return
	}

	server, err := selectServerRecordById(serverId)

	if err != nil {
		httputils.RespondWithInternalServerError(context, err)
		return
	}

	err = populateServerWithProperties(server)

	if err != nil {
		httputils.RespondWithInternalServerError(context, err)
		return
	}

	err = populateServerWithEulaAcceptanceStatus(server)

	if err != nil {
		httputils.RespondWithInternalServerError(context, err)
		return
	}

	httputils.RespondWithStatusOk(context, server)
}

func selectServerRecordsByUserId(id string) (*[]MCServerLite, error) {
	db, err := sql.Open("sqlite3", "./gomine.db")

	if err != nil {
		return nil, err
	}

	defer db.Close()

	statement, err := db.Prepare("select name, runtime, path, pid, status, user_id from servers where user_id=?")

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	rows, err := statement.Query(id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	servers := []MCServerLite{}

	for rows.Next() {
		server := MCServerLite{ID: id}
		err = rows.Scan(&server.Name, &server.Runtime, &server.Path, &server.PID, &server.Status, &server.UserID)

		if err == nil {
			servers = append(servers, server)
		}
	}

	return &servers, nil
}

func GetServersByUserId(context *gin.Context) {
	userId := context.Query("u")

	if userId == "" {
		httputils.RespondWithNotFound(context, errors.New("GetServersByUserId: No user id provided."))
		return
	}

	servers, err := selectServerRecordsByUserId(userId)

	if err != nil {
		httputils.RespondWithInternalServerError(context, err)
		return
	}

	httputils.RespondWithStatusOk(context, servers)
}
